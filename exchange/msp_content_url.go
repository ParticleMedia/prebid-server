package exchange

import (
	"github.com/buger/jsonparser"
)

const (
	contentURLPresent = "present"
	contentURLNone    = "none"
)

// contentURLPrefix reports whether an outgoing bidder request body carries an app (or site) content URL.
// It returns one of two fixed values so the label cardinality can't grow with the URLs.
func contentURLPrefix(body []byte) string {
	url, err := jsonparser.GetString(body, "app", "content", "url")
	if err != nil || url == "" {
		url, err = jsonparser.GetString(body, "site", "content", "url")
	}
	if err != nil || url == "" {
		return contentURLNone
	}
	return contentURLPresent
}
