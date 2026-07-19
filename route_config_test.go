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
		ApplyKeywordRules:    true,
		FollowAuthorsEnabled: true,
		FollowedAuthors:      []string{"  张三,李四 ", "李四"},
	}}
	got := req.applyTo(file.Config{})
	if got.GlobalHot.WindowHours != 3 {
		t.Fatalf("WindowHours = %d, want default 3 for unsupported value", got.GlobalHot.WindowHours)
	}
	if got.GlobalHot.MinCommentNum != 100 {
		t.Fatalf("MinCommentNum = %d, want 100 for an in-range lower preset", got.GlobalHot.MinCommentNum)
	}
	if len(got.GlobalHot.FollowedAuthors) != 2 || got.GlobalHot.FollowedAuthors[0] != "张三" || got.GlobalHot.FollowedAuthors[1] != "李四" {
		t.Fatalf("FollowedAuthors = %#v", got.GlobalHot.FollowedAuthors)
	}

	response := productConfigFromConfig(got)
	if !response.GlobalHot.Enabled || !response.GlobalHot.ApplyKeywordRules || !response.GlobalHot.FollowAuthorsEnabled {
		t.Fatal("global hot flags were not preserved")
	}
}
