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
		"id=\"globalHotWindow\" type=\"number\"",
		"id=\"globalHotMinComment\" type=\"number\"",
		"id=\"hotKeywordTokenBox\"",
		"id=\"followedAuthorsTokenBox\"",
		"id=\"authorKeywordTokenBox\"",
		"id=\"productKeywordField\"",
		"id=\"discoveryTimeEditor\"",
		"id=\"discoveryFieldsRow\"",
		"discovery-fields-row",
		"repeat(auto-fit,minmax(190px,1fr))",
		"function bindDiscoveryTokenEditor",
		"function normalizeGlobalHot(value)",
		"var selectedRuleKind = 'product'",
		"function openRuleTypePicker",
		"offset: [top + 'px', left + 'px']",
		"data-rule-kind=\"hot\"",
		"id=\"hotRuleEditor\"",
		"id=\"authorRuleEditor\"",
		"/discoverySearch",
	}
	for _, marker := range required {
		if !strings.Contains(page, marker) {
			t.Errorf("dashboard HTML missing search feedback marker %q", marker)
		}
	}
	if strings.Contains(page, "本次临时搜索") {
		t.Error("dashboard still shows the obsolete label 本次临时搜索")
	}
	if strings.Contains(page, "id=\"globalHotEnabled\"") || strings.Contains(page, "id=\"followAuthorsEnabled\"") {
		t.Error("discovery editors still contain redundant enable switches")
	}
	if strings.Contains(page, "product-author") {
		t.Error("search result cards should not display author metadata")
	}
	for _, marker := range []string{"<strong>搜索关键词</strong>", "<strong>搜索热门</strong>", "<strong>搜索作者</strong>", "for=\"followedAuthorsInput\">搜索作者</label>"} {
		if !strings.Contains(page, marker) {
			t.Errorf("dashboard HTML missing discovery search label %q", marker)
		}
	}
	if strings.Contains(page, "<select id=\"globalHotWindow\"") || strings.Contains(page, "<select id=\"globalHotMinComment\"") {
		t.Error("global hot numeric settings still use preset selects")
	}
}
