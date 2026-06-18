package push

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"

	"ggball.com/smzdm/file"
	"ggball.com/smzdm/smzdm"
)

var telegramEndpointBase = "https://api.telegram.org"
var httpClient = http.DefaultClient

func PushProducts(pro []smzdm.Product, conf file.Config) {
	PushProWithTelegram(pro, conf)
}

func PushText(resText string, conf file.Config) {
	PushTextWithTelegram(resText, conf)
}

func PushTargetProducts(pro []smzdm.Product, conf file.Config, atMobiles []string) {
	PushTargetWithTelegram(pro, conf)
}

func PushProWithTelegram(pro []smzdm.Product, conf file.Config) {
	if len(pro) == 0 || !canPushTelegram(conf) {
		return
	}
	limit := productLimit(len(pro), conf.SatisfyNum)
	content := buildProductsHTML(pro[:limit], "什么值得买好价")
	pushTelegram(content, conf)
}

func PushTargetWithTelegram(pro []smzdm.Product, conf file.Config) {
	if len(pro) == 0 || !canPushTelegram(conf) {
		return
	}
	content := buildProductsHTML(pro, "好物到了")
	pushTelegram(content, conf)
}

func PushTextWithTelegram(resText string, conf file.Config) {
	if strings.TrimSpace(resText) == "" || !canPushTelegram(conf) {
		return
	}
	pushTelegram(html.EscapeString(resText)+"\n\n什么值得买", conf)
}

func pushTelegram(content string, conf file.Config) {
	params := TelegramMessageParam{
		ChatID:                conf.Telegram.ChatID,
		Text:                  content,
		ParseMode:             telegramParseMode(conf.Telegram.ParseMode),
		DisableWebPagePreview: conf.Telegram.DisableWebPagePreview,
	}

	paramsJson, _ := json.Marshal(params)
	endpoint := strings.TrimRight(telegramEndpointBase, "/") + "/bot" + strings.TrimSpace(conf.Telegram.BotToken) + "/sendMessage"
	resp, err := httpClient.Post(endpoint, "application/json;charset=utf-8", bytes.NewBuffer(paramsJson))
	if err != nil {
		fmt.Println("Telegram push failed:", err)
		return
	}
	defer resp.Body.Close()

	contentBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Telegram response read failed:", err)
		return
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println("Telegram response status:", resp.StatusCode, string(contentBytes))
		return
	}
	fmt.Println(string(contentBytes))
}

func canPushTelegram(conf file.Config) bool {
	return conf.Telegram.Enabled &&
		strings.TrimSpace(conf.Telegram.BotToken) != "" &&
		strings.TrimSpace(conf.Telegram.ChatID) != ""
}

func buildProductsHTML(pro []smzdm.Product, title string) string {
	var builder strings.Builder
	builder.WriteString("<b>" + html.EscapeString(title) + "</b>\n")
	for _, item := range pro {
		builder.WriteString("\n")
		builder.WriteString(formatProductHTML(item))
	}
	return builder.String()
}

func formatProductHTML(item smzdm.Product) string {
	title := html.EscapeString(item.ArticleTitle)
	link := html.EscapeString(item.ArticleUrl)
	price := html.EscapeString(item.ArticlePrice)
	comment := html.EscapeString(item.ArticleComment)
	worthy := html.EscapeString(item.ArticleWorthy)
	return fmt.Sprintf("<a href=\"%s\">%s</a>\n%s, 评论 %s, 值率 %s%%\n", link, title, price, comment, worthy)
}

func telegramParseMode(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "HTML"
	}
	return value
}

func productLimit(length int, satisfyNum int) int {
	if satisfyNum <= 0 || satisfyNum > length {
		return length
	}
	return satisfyNum
}
