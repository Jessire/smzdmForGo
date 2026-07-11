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
		"<strong>本次搜索</strong>",
		"toggleClass('is-searching', state === '搜索中')",
		"if (selector === '#searchRule')",
	}
	for _, marker := range required {
		if !strings.Contains(page, marker) {
			t.Errorf("dashboard HTML missing search feedback marker %q", marker)
		}
	}
	if strings.Contains(page, "本次临时搜索") {
		t.Error("dashboard still shows the obsolete label 本次临时搜索")
	}
}
