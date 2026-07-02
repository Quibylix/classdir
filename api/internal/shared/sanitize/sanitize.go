package sanitize

import (
	"github.com/microcosm-cc/bluemonday"
	"regexp"
)

var RevealPolicy *bluemonday.Policy

func init() {
	RevealPolicy = bluemonday.UGCPolicy()

	RevealPolicy.AllowAttrs("class").Globally()

	RevealPolicy.AllowAttrs(
		"data-transition", "data-autoplay", "data-fragment-index", "data-id",
	).Globally()

	urlRegex := regexp.MustCompile(`(?i)^(https?://|/|#).*`)

	RevealPolicy.AllowAttrs("data-background-image", "data-background-video").Matching(urlRegex).Globally()

	colorRegex := regexp.MustCompile(`(?i)^(#([0-9a-f]{3,4}|[0-9a-f]{6}|[0-9a-f]{8})|rgb.*|hsl.*|[a-z]+)$`)
	RevealPolicy.AllowAttrs("data-background-color").Matching(colorRegex).Globally()
}
