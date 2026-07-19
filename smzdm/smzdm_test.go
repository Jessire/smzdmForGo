package smzdm

import (
	"fmt"
	"testing"
	"time"

	"ggball.com/smzdm/file"
)

func TestGlobalHotFiltersAndAuthorMatching(t *testing.T) {
	enabled := true
	globalConf = file.Config{KeywordRules: []file.KeywordRule{{Enabled: &enabled, Words: []string{"显示器"}}}}
	conf := file.Config{GlobalHot: file.GlobalHotConfig{
		Enabled:              true,
		MinCommentNum:        200,
		ApplyKeywordRules:    true,
		FollowAuthorsEnabled: true,
		FollowedAuthors:      []string{"  关注的人  "},
	}}
	hot := Product{ArticleTitle: "27寸显示器", ArticleComment: "200", ArticleDate: "1784434480"}
	low := Product{ArticleTitle: "27寸显示器", ArticleComment: "199", ArticleDate: "1784434480"}
	author := Product{ArticleTitle: "完全不相关", ArticleComment: "0", Referral: "关注的人", ArticleDate: "1784434480"}
	if !globalFeedProductMatches(hot, conf) {
		t.Fatal("expected 200-comment rule match")
	}
	if globalFeedProductMatches(low, conf) {
		t.Fatal("expected item below comment floor to be rejected")
	}
	if !globalFeedProductMatches(author, conf) {
		t.Fatal("expected followed author item to bypass hot floor")
	}
	if !followedAuthorMatch("关注的人", []string{" 关注的人 "}) {
		t.Fatal("expected normalized author match")
	}
	globalConf = file.Config{}
	if !globalFeedProductMatches(hot, conf) {
		t.Fatal("expected hot item to pass when no secondary rules exist")
	}
}

func TestSortProductsByCommentAndTime(t *testing.T) {
	later := time.Now().Add(-time.Minute).Unix()
	earlier := time.Now().Add(-time.Hour).Unix()
	items := []Product{
		{ArticleId: "old", ArticleComment: "200", ArticleDate: stringInt64(earlier)},
		{ArticleId: "new", ArticleComment: "200", ArticleDate: stringInt64(later)},
		{ArticleId: "top", ArticleComment: "201", ArticleDate: stringInt64(earlier)},
	}
	sortProductsByComment(items)
	if items[0].ArticleId != "top" || items[1].ArticleId != "new" {
		t.Fatalf("unexpected sorted order: %#v", items)
	}
}

func stringInt64(value int64) string {
	return fmt.Sprintf("%d", value)
}

func TestParseMetric(t *testing.T) {
	tests := map[string]int{
		"12":   12,
		"1.2k": 1200,
		"3万":   30000,
		"68%":  68,
		"暂无":   0,
	}

	for input, expected := range tests {
		if got := parseMetric(input); got != expected {
			t.Fatalf("parseMetric(%q)=%d, want %d", input, got, expected)
		}
	}
}

func TestPriceInRange(t *testing.T) {
	good := Product{ArticlePrice: "到手价 1,299.50 元"}
	if !priceInRange(good, 1000, 1300) {
		t.Fatal("expected price to be in range")
	}
	if priceInRange(good, 0, 1000) {
		t.Fatal("expected price to be above max range")
	}
}

func TestKeywordRules(t *testing.T) {
	enabled := true
	lowComment := 5
	lowWorthy := 20
	minPrice := 300.0
	maxPrice := 2000.0
	globalConf = file.Config{
		KeyWords:      []string{"显示器"},
		LowCommentNum: 1,
		LowWorthyNum:  6,
		KeywordRules: []file.KeywordRule{
			{
				Enabled:       &enabled,
				Words:         []string{"显示器"},
				FilterWords:   []string{"二手"},
				LowCommentNum: &lowComment,
				LowWorthyNum:  &lowWorthy,
				MinPrice:      &minPrice,
				MaxPrice:      &maxPrice,
			},
		},
	}

	good := Product{
		ArticleTitle:   "27寸 4K 显示器",
		ArticlePrice:   "1299元",
		ArticleComment: "12",
		ArticleWorthy:  "88",
	}
	if !matchesPersonalRules(good) {
		t.Fatal("expected product to match keyword rule")
	}

	good.ArticleTitle = "二手 27寸 4K 显示器"
	if matchesPersonalRules(good) {
		t.Fatal("expected product with rule filter word to be rejected")
	}

	enabled = false
	good.ArticleTitle = "27寸 4K 显示器"
	if matchesPersonalRules(good) {
		t.Fatal("expected disabled keyword rule to be rejected")
	}
}
