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
}
