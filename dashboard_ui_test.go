package main

import (
	"os"
	"strings"
	"testing"
)

// The dashboard was rebuilt (2026-07) as a vanilla-JS three-column layout.
// These markers assert the new UI's structural invariants.
func TestDashboardStructure(t *testing.T) {
	html, err := os.ReadFile("template/html/index.html")
	if err != nil {
		t.Fatalf("read dashboard HTML: %v", err)
	}
	page := string(html)

	required := []string{
		// layout containers
		"id=\"ruleList\"",
		"id=\"ruleEditor\"",
		"id=\"searchRule\"",
		"id=\"tgPanel\"",
		"id=\"statsBar\"",
		// rule editors per kind
		"id=\"productKeywordField\"",
		"id=\"hotRuleEditor\"",
		"id=\"authorRuleEditor\"",
		"id=\"discoveryTimeEditor\"",
		"id=\"globalHotWindow\" type=\"number\"",
		"id=\"globalHotMinComment\" type=\"number\"",
		"id=\"hotKeywordTokenBox\"",
		"id=\"followedAuthorsTokenBox\"",
		"id=\"authorKeywordTokenBox\"",
		"data-rule-kind",
		"var selectedRuleKind = 'product'",
		// discovery labels
		"搜索热门",
		"搜索作者",
		"for=\"followedAuthorsInput\">搜索作者</label>",
		// backend endpoints the page must call
		"/productConfig",
		"/productSearch",
		"/discoverySearch",
		"/telegramTest",
		"/pushLogs",
		"/imageProxy",
		// search feedback
		"rule-auto-state",
		"is-searching",
		// theming
		"data-theme=\"dark\"",
	}
	for _, marker := range required {
		if !strings.Contains(page, marker) {
			t.Errorf("dashboard HTML missing marker %q", marker)
		}
	}

	// obsolete pieces that must not come back
	if strings.Contains(page, "本次临时搜索") {
		t.Error("dashboard still shows the obsolete label 本次临时搜索")
	}
	if strings.Contains(page, "$.ajax") || strings.Contains(page, "jquery") {
		t.Error("dashboard should not depend on jQuery")
	}
	if strings.Contains(page, "id=\"globalHotEnabled\"") || strings.Contains(page, "id=\"followAuthorsEnabled\"") {
		t.Error("discovery editors must not contain redundant enable switches")
	}
	if strings.Contains(page, "<select id=\"globalHotWindow\"") || strings.Contains(page, "<select id=\"globalHotMinComment\"") {
		t.Error("global hot numeric settings must be number inputs, not preset selects")
	}
	if strings.Contains(page, "product-author") {
		t.Error("search result cards should not display author metadata")
	}
}
