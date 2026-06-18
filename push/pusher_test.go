package push

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ggball.com/smzdm/file"
	"ggball.com/smzdm/smzdm"
)

func TestPushProWithTelegram(t *testing.T) {
	var got TelegramMessageParam
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method=%s, want POST", r.Method)
		}
		gotPath = r.URL.Path
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	oldEndpoint := telegramEndpointBase
	telegramEndpointBase = server.URL
	defer func() { telegramEndpointBase = oldEndpoint }()

	conf := file.Config{
		SatisfyNum: 1,
		Telegram: file.Telegram{
			Enabled:               true,
			BotToken:              "123:test",
			ChatID:                "456",
			ParseMode:             "HTML",
			DisableWebPagePreview: true,
		},
	}
	PushProWithTelegram([]smzdm.Product{
		{
			ArticleTitle:   "测试商品",
			ArticlePrice:   "99元",
			ArticleComment: "10",
			ArticleWorthy:  "80",
			ArticlePic:     "https://example.com/a.jpg",
			ArticleUrl:     "https://example.com/a",
		},
	}, conf)

	if gotPath != "/bot123:test/sendMessage" {
		t.Fatalf("path=%s, want Telegram sendMessage path", gotPath)
	}
	if got.ChatID != "456" {
		t.Fatalf("chat_id=%s, want 456", got.ChatID)
	}
	if got.ParseMode != "HTML" {
		t.Fatalf("parse_mode=%s, want HTML", got.ParseMode)
	}
	if !got.DisableWebPagePreview {
		t.Fatal("disable_web_page_preview=false, want true")
	}
	if !strings.Contains(got.Text, "<a href=\"https://example.com/a\">测试商品</a>") {
		t.Fatalf("text=%s, want product link", got.Text)
	}
	if !strings.Contains(got.Text, "99元, 评论 10, 值率 80%") {
		t.Fatalf("text=%s, want product metrics", got.Text)
	}
}
