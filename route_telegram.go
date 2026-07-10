package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"ggball.com/smzdm/file"
)

var telegramAvatarEndpointBase = "https://api.telegram.org"
var telegramAvatarHTTPClient = &http.Client{Timeout: 10 * time.Second}

type telegramChatResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
	Result      struct {
		Photo *struct {
			SmallFileID string `json:"small_file_id"`
		} `json:"photo"`
	} `json:"result"`
}

type telegramFileResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
	Result      struct {
		FilePath string `json:"file_path"`
	} `json:"result"`
}

func TelegramAvatarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	avatar, contentType, err := fetchTelegramAvatar(r.Context(), currentConfig().Telegram)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "private, max-age=300")
	w.Write(avatar)
}

func fetchTelegramAvatar(ctx context.Context, telegram file.Telegram) ([]byte, string, error) {
	token := strings.TrimSpace(telegram.BotToken)
	chatID := strings.TrimSpace(telegram.ChatID)
	if token == "" || chatID == "" {
		return nil, "", fmt.Errorf("telegram 未配置")
	}

	var chat telegramChatResponse
	chatURL := strings.TrimRight(telegramAvatarEndpointBase, "/") + "/bot" + token + "/getChat?chat_id=" + url.QueryEscape(chatID)
	if err := getTelegramJSON(ctx, chatURL, &chat); err != nil {
		return nil, "", err
	}
	if !chat.OK || chat.Result.Photo == nil || strings.TrimSpace(chat.Result.Photo.SmallFileID) == "" {
		return nil, "", fmt.Errorf("telegram 账户没有可用头像")
	}

	var fileInfo telegramFileResponse
	fileURL := strings.TrimRight(telegramAvatarEndpointBase, "/") + "/bot" + token + "/getFile?file_id=" + url.QueryEscape(chat.Result.Photo.SmallFileID)
	if err := getTelegramJSON(ctx, fileURL, &fileInfo); err != nil {
		return nil, "", err
	}
	if !fileInfo.OK || strings.TrimSpace(fileInfo.Result.FilePath) == "" {
		return nil, "", fmt.Errorf("telegram 头像文件不可用")
	}

	avatarURL := strings.TrimRight(telegramAvatarEndpointBase, "/") + "/file/bot" + token + "/" + strings.TrimLeft(fileInfo.Result.FilePath, "/")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, avatarURL, nil)
	if err != nil {
		return nil, "", err
	}
	resp, err := telegramAvatarHTTPClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", fmt.Errorf("telegram 头像响应状态 %d", resp.StatusCode)
	}
	avatar, err := io.ReadAll(io.LimitReader(resp.Body, 5<<20))
	if err != nil {
		return nil, "", err
	}
	if len(avatar) == 0 {
		return nil, "", fmt.Errorf("telegram 头像为空")
	}
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(avatar)
	}
	return avatar, contentType, nil
}

func getTelegramJSON(ctx context.Context, endpoint string, target interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	resp, err := telegramAvatarHTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram api 响应状态 %d", resp.StatusCode)
	}
	return json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(target)
}
