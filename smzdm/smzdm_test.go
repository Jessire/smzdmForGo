package smzdm

import (
	"testing"

	"ggball.com/smzdm/file"
)

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

func TestContainsWordRegex(t *testing.T) {
	good := Product{ArticleTitle: "27寸 4K 显示器 IPS 面板"}

	if !containsWord(good, "显示器") {
		t.Fatal("expected plain keyword match")
	}
	if !containsWord(good, "显示器|屏幕") {
		t.Fatal("expected regex OR to match 显示器")
	}
	if !containsWord(good, `4K|OLED`) {
		t.Fatal("expected regex OR to match 4K")
	}
	if containsWord(good, `OLED|MiniLED`) {
		t.Fatal("expected non-matching regex OR to fail")
	}
	if !containsWord(good, `(?i)ips`) {
		t.Fatal("expected case-insensitive regex match")
	}

	// Invalid regex should fall back to literal substring
	if !containsWord(Product{ArticleTitle: "price (unclosed"}, `(unclosed`) {
		t.Fatal("expected invalid regex to fall back to literal contains")
	}
}

func TestContainsWordRegexFilter(t *testing.T) {
	enabled := true
	globalConf = file.Config{
		KeywordRules: []file.KeywordRule{
			{
				Enabled:     &enabled,
				Words:       []string{`显示器|屏幕`},
				FilterWords: []string{`二手|翻新`},
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
		t.Fatal("expected regex keyword to match")
	}

	good.ArticleTitle = "二手 27寸 4K 显示器"
	if matchesPersonalRules(good) {
		t.Fatal("expected regex filter word to reject")
	}

	good.ArticleTitle = "屏幕 挂灯"
	if !matchesPersonalRules(good) {
		t.Fatal("expected alternate regex branch 屏幕 to match")
	}
}

func TestSearchSeeds(t *testing.T) {
	if got := searchSeeds("显示器"); len(got) != 1 || got[0] != "显示器" {
		t.Fatalf("plain seed = %#v, want [显示器]", got)
	}
	got := searchSeeds("显示器|屏幕")
	if len(got) != 2 || got[0] != "显示器" || got[1] != "屏幕" {
		t.Fatalf("OR seeds = %#v, want [显示器 屏幕]", got)
	}
	got = searchSeeds(`(4K|8K).{0,6}显示器`)
	if len(got) == 0 || got[0] == "" {
		t.Fatalf("complex seed empty: %#v", got)
	}
	// Longest literal should be preferred from first alternative group content
	if !containsWord(Product{ArticleTitle: "4K 显示器"}, `(4K|8K).{0,6}显示器`) {
		t.Fatal("expected complex regex to match title")
	}
}

func TestSearchProductMatchesKeyword(t *testing.T) {
	rule := file.KeywordRule{
		Words:       []string{`显示器|屏幕`},
		FilterWords: nil,
	}
	if !searchProductMatches(Product{ArticleTitle: "曲面屏幕 144Hz"}, rule) {
		t.Fatal("expected local regex keyword filter to pass")
	}
	if searchProductMatches(Product{ArticleTitle: "无线鼠标"}, rule) {
		t.Fatal("expected unrelated title to fail keyword filter")
	}
}
