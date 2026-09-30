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
		{"app content url", `{"app":{"content":{"url":"https://www.example.com/a"}}}`, "present"},
		{"site content url", `{"site":{"content":{"url":"http://example.com"}}}`, "present"},
		{"empty app url falls back to site", `{"app":{"content":{"url":""}},"site":{"content":{"url":"http://site.com"}}}`, "present"},
		{"leading whitespace", " \n{\"app\":{\"content\":{\"url\":\"https://x.com\"}}}", "present"},
		{"non-url value", `{"app":{"content":{"url":"日本語\u0000abc"}}}`, "present"},
		{"empty url", `{"app":{"content":{"url":""}}}`, "none"},
		{"no content", `{"app":{"bundle":"com.example"}}`, "none"},
		{"url not a string", `{"app":{"content":{"url":5}}}`, "none"},
		{"no body", "", "none"},
		{"not json", "\x12\x04test", "none"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, contentURLPrefix([]byte(tt.body)))
		})
	}
}
