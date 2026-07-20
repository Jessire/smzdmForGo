package main

import (
	"testing"

	"ggball.com/smzdm/file"
)

func TestGlobalHotConfigNormalization(t *testing.T) {
	req := productConfigRequest{GlobalHot: globalHotConfigRequest{
		Enabled:              true,
		WindowHours:          9,
		MinCommentNum:        175,
		HotKeywords:          []string{" 显示器, 面包 "},
		FollowAuthorsEnabled: true,
		FollowedAuthors:      []string{"  张三,李四 ", "李四"},
		AuthorKeywords:       []string{" 耳机 "},
		FilterWords:          []string{" 临期, 售罄 "},
		LowCommentNum:        3,
		LowWorthyNum:         8,
		MinPrice:             29.9,
		MaxPrice:             199,
	}}
	got := req.applyTo(file.Config{})
	if got.GlobalHot.WindowHours != 9 {
		t.Fatalf("WindowHours = %d, want user-entered 9", got.GlobalHot.WindowHours)
	}
	if got.GlobalHot.MinCommentNum != 175 {
		t.Fatalf("MinCommentNum = %d, want user-entered 175", got.GlobalHot.MinCommentNum)
	}
	if len(got.GlobalHot.FollowedAuthors) != 2 || got.GlobalHot.FollowedAuthors[0] != "张三" || got.GlobalHot.FollowedAuthors[1] != "李四" {
		t.Fatalf("FollowedAuthors = %#v", got.GlobalHot.FollowedAuthors)
	}

	response := productConfigFromConfig(got)
	if !response.GlobalHot.Enabled || !response.GlobalHot.FollowAuthorsEnabled {
		t.Fatal("global hot flags were not preserved")
	}
	if len(response.GlobalHot.HotKeywords) != 2 || len(response.GlobalHot.AuthorKeywords) != 1 {
		t.Fatalf("independent keywords were not preserved: %#v", response.GlobalHot)
	}
	if len(response.GlobalHot.FilterWords) != 2 || response.GlobalHot.LowCommentNum != 3 || response.GlobalHot.LowWorthyNum != 8 || response.GlobalHot.MinPrice != 29.9 || response.GlobalHot.MaxPrice != 199 {
		t.Fatalf("shared discovery filters were not preserved: %#v", response.GlobalHot)
	}
}

func TestGlobalHotConfigUsesDefaultsForNonPositiveValues(t *testing.T) {
	req := productConfigRequest{GlobalHot: globalHotConfigRequest{
		WindowHours:   -6,
		MinCommentNum: 0,
	}}
	got := req.applyTo(file.Config{})
	if got.GlobalHot.WindowHours != 3 {
		t.Fatalf("WindowHours = %d, want default 3", got.GlobalHot.WindowHours)
	}
	if got.GlobalHot.MinCommentNum != 200 {
		t.Fatalf("MinCommentNum = %d, want default 200", got.GlobalHot.MinCommentNum)
	}
}
