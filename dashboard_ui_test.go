package main

import (
	"os"
	"strings"
	"testing"
)

func TestDashboardSearchFeedbackPlacement(t *testing.T) {
	html, err := os.ReadFile("template/html/index.html")
	if err != nil {
		t.Fatalf("read dashboard HTML: %v", err)
	}
	page := string(html)
	required := []string{
		"grid-auto-rows:112px",
		".product-card{min-height:112px;height:112px}",
		".rule-auto-state.is-searching:before",
		"toggleClass('is-searching', state === '搜索中')",
		"if (selector === '#searchRule')",
		"id=\"globalHotEnabled\"",
		"id=\"globalHotWindow\" type=\"number\"",
		"id=\"globalHotMinComment\" type=\"number\"",
		"id=\"hotKeywordTokenBox\"",
		"id=\"followedAuthorsTokenBox\"",
		"id=\"followAuthorsEnabled\"",
		"id=\"authorKeywordTokenBox\"",
		"function bindDiscoveryTokenEditor",
		"function normalizeGlobalHot(value)",
	}
	for _, marker := range required {
		if !strings.Contains(page, marker) {
			t.Errorf("dashboard HTML missing search feedback marker %q", marker)
		}
	}
	if strings.Contains(page, "本次临时搜索") {
		t.Error("dashboard still shows the obsolete label 本次临时搜索")
	}
	if strings.Contains(page, "<select id=\"globalHotWindow\"") || strings.Contains(page, "<select id=\"globalHotMinComment\"") {
		t.Error("global hot numeric settings still use preset selects")
	}
}
