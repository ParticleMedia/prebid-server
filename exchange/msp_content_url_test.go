package exchange

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContentURLPrefix(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{"app content url", `{"app":{"content":{"url":"https://www.example.com/a"}}}`, "https:"},
		{"site content url", `{"site":{"content":{"url":"http://example.com"}}}`, "http:/"},
		{"app wins over site", `{"site":{"content":{"url":"http://site"}},"app":{"content":{"url":"https://app"}}}`, "https:"},
		{"leading whitespace", " \n{\"app\":{\"content\":{\"url\":\"https://x\"}}}", "https:"},
		{"shorter than prefix", `{"app":{"content":{"url":"abc"}}}`, "abc"},
		{"multibyte runes", `{"app":{"content":{"url":"日本語のニュース"}}}`, "日本語のニュ"},
		{"invalid utf8", "{\"app\":{\"content\":{\"url\":\"\xffabcdefg\"}}}", "�abcde"},
		{"empty url", `{"app":{"content":{"url":""}}}`, "none"},
		{"no content", `{"app":{"bundle":"com.example"}}`, "none"},
		{"url not a string", `{"app":{"content":{"url":5}}}`, "none"},
		{"no body", "", "no_body"},
		{"not json", "\x12\x04test", "not_json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, contentURLPrefix([]byte(tt.body)))
		})
	}
}
