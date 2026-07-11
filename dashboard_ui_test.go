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
		".editor-panel .editor-actions #searchRule{flex:0 0 min(95px,12%);min-width:95px;padding-inline:4px;font-size:13px;gap:5px}",
		".editor-panel .editor-actions #saveProductConfig{flex:0 0 190px;min-width:190px}",
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
