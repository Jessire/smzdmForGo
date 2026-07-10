package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"ggball.com/smzdm/file"
)

func TestFetchTelegramAvatarProxiesConfiguredChatPhoto(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/bot123:test/getChat":
			if r.URL.Query().Get("chat_id") != "456" {
				t.Fatalf("chat_id=%q, want 456", r.URL.Query().Get("chat_id"))
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"ok":true,"result":{"photo":{"small_file_id":"avatar-file"}}}`))
		case "/bot123:test/getFile":
			if r.URL.Query().Get("file_id") != "avatar-file" {
				t.Fatalf("file_id=%q, want avatar-file", r.URL.Query().Get("file_id"))
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"ok":true,"result":{"file_path":"photos/avatar.jpg"}}`))
		case "/file/bot123:test/photos/avatar.jpg":
			w.Header().Set("Content-Type", "image/jpeg")
			w.Write([]byte("avatar-bytes"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	oldBase := telegramAvatarEndpointBase
	oldClient := telegramAvatarHTTPClient
	telegramAvatarEndpointBase = server.URL
	telegramAvatarHTTPClient = server.Client()
	defer func() {
		telegramAvatarEndpointBase = oldBase
		telegramAvatarHTTPClient = oldClient
	}()

	avatar, contentType, err := fetchTelegramAvatar(context.Background(), file.Telegram{
		BotToken: "123:test",
		ChatID:   "456",
	})
	if err != nil {
		t.Fatalf("fetchTelegramAvatar() error = %v", err)
	}
	if string(avatar) != "avatar-bytes" {
		t.Fatalf("avatar=%q, want avatar-bytes", string(avatar))
	}
	if contentType != "image/jpeg" {
		t.Fatalf("contentType=%q, want image/jpeg", contentType)
	}
}
