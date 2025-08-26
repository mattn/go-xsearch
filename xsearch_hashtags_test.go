package xsearch

import (
	"reflect"
	"testing"
)

func TestExtractHashtags(t *testing.T) {
	entries := []Entry{
		{Hashtags: []any{"go", "golang", "dev"}},
		{Hashtags: []any{"go", "code"}},
		{Hashtags: []any{"golang"}},
	}
	got := ExtractHashtags(entries)
	expected := []string{"go", "golang", "dev", "code"}
	// Convert to map for comparison (order doesn't matter)
	gotMap := make(map[string]struct{})
	expMap := make(map[string]struct{})
	for _, v := range got {
		gotMap[v] = struct{}{}
	}
	for _, v := range expected {
		expMap[v] = struct{}{}
	}
	if !reflect.DeepEqual(gotMap, expMap) {
		t.Errorf("ExtractHashtags() = %v, want %v", got, expected)
	}
}
