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

func TestPushProWithTelegramBatches(t *testing.T) {
	var texts []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var msg TelegramMessageParam
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		texts = append(texts, msg.Text)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	oldEndpoint := telegramEndpointBase
	telegramEndpointBase = server.URL
	defer func() { telegramEndpointBase = oldEndpoint }()

	// Clear recent logs for a clean assertion.
	logMu.Lock()
	recentLogs = nil
	logMu.Unlock()

	products := make([]smzdm.Product, 0, 7)
	for i := 1; i <= 7; i++ {
		products = append(products, smzdm.Product{
			ArticleTitle:   "商品" + string(rune('0'+i)),
			ArticlePrice:   "10元",
			ArticleComment: "1",
			ArticleWorthy:  "10",
			ArticleUrl:     "https://example.com/p" + string(rune('0'+i)),
		})
	}
	conf := file.Config{
		SatisfyNum: 3,
		Telegram: file.Telegram{
			Enabled:  true,
			BotToken: "123:test",
			ChatID:   "456",
		},
	}
	PushProWithTelegram(products, conf)

	if len(texts) != 3 {
		t.Fatalf("messages=%d, want 3 batches (3+3+1)", len(texts))
	}
	if !strings.Contains(texts[0], "1/3") || !strings.Contains(texts[1], "2/3") || !strings.Contains(texts[2], "3/3") {
		t.Fatalf("batch titles missing: %#v", texts)
	}
	logs := RecentLogs()
	if len(logs) < 7 {
		t.Fatalf("logs=%d, want at least 7 product entries", len(logs))
	}
}
