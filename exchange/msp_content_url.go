package exchange

import (
	"bytes"
	"strings"

	"github.com/buger/jsonparser"
)

const contentURLPrefixLen = 11

// contentURLPrefix returns the first characters of the app (or site) content URL in an outgoing bidder request body.
func contentURLPrefix(body []byte) string {
	body = bytes.TrimLeft(body, " \t\r\n")
	if len(body) == 0 {
		return "no_body"
	}
	if body[0] != '{' {
		return "not_json"
	}

	url, err := jsonparser.GetString(body, "app", "content", "url")
	if err != nil {
		url, err = jsonparser.GetString(body, "site", "content", "url")
	}
	if err != nil || url == "" {
		return "none"
	}

	n := 0
	for i := range url {
		if n == contentURLPrefixLen {
			url = url[:i]
			break
		}
		n++
	}
	return strings.ToValidUTF8(url, "�")
}
