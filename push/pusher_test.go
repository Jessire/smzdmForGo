package push

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ggball.com/smzdm/file"
	"ggball.com/smzdm/smzdm"
)

func TestPushProWithWxPusher(t *testing.T) {
	var got WxPusherMessageParam
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method=%s, want POST", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Write([]byte(`{"code":1000,"msg":"ok"}`))
	}))
	defer server.Close()

	oldEndpoint := wxPusherEndpoint
	wxPusherEndpoint = server.URL
	defer func() { wxPusherEndpoint = oldEndpoint }()

	conf := file.Config{
		SatisfyNum: 1,
		WxPusher: file.WxPusher{
			Enabled:     true,
			AppToken:    "AT_test",
			UIDs:        []string{"UID_test"},
			ContentType: 3,
		},
	}
	PushProWithWxPusher([]smzdm.Product{
		{
			ArticleTitle:   "测试商品",
			ArticlePrice:   "99元",
			ArticleComment: "10",
			ArticleWorthy:  "80",
			ArticlePic:     "https://example.com/a.jpg",
			ArticleUrl:     "https://example.com/a",
		},
	}, conf)

	if got.AppToken != "AT_test" {
		t.Fatalf("appToken=%s, want AT_test", got.AppToken)
	}
	if len(got.UIDs) != 1 || got.UIDs[0] != "UID_test" {
		t.Fatalf("uids=%v, want UID_test", got.UIDs)
	}
	if got.ContentType != 3 {
		t.Fatalf("contentType=%d, want 3", got.ContentType)
	}
	if got.URL != "https://example.com/a" {
		t.Fatalf("url=%s, want product url", got.URL)
	}
}
