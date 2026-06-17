package file

import (
	"path/filepath"
	"testing"
)

func TestReadAndWritePushed(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pushed.json")

	readPushedInfo := ReadPusedInfo(path)
	if len(readPushedInfo) != 0 {
		t.Fatalf("new pushed map len=%d, want 0", len(readPushedInfo))
	}

	pushedMap := map[string]interface{}{
		"222": "222",
	}
	WritePushedInfo(pushedMap, readPushedInfo, path)

	got := ReadPusedInfo(path)
	if got["222"] != "222" {
		t.Fatalf("pushed value=%v, want 222", got["222"])
	}
}
