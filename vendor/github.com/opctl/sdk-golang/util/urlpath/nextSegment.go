// Package urlpath exports NextSegment for nexting urlpath's.
// This negates the need for global &/or sub HTTP Muxers (which generally have hidden coupling to handlers)
// Equiped w/ urlpath.NextSegment, handlers can simply next the path of their req url, do processing specific to them,
// and delegate as needed to any children, which follow suit.
// this idea was inspired by https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html
package urlpath

import (
	"net/url"
	"strings"
)

// NextSegment mutates the provided URL by consuming the next segment of it's EscapedPath.
// returns the consumed segment, in unescaped form
func NextSegment(
	u *url.URL,
) (string, error) {
	escapedPathParts := strings.SplitN(u.EscapedPath(), "/", 2)
	if len(escapedPathParts) == 1 {
		// this was the last path segment of path

		// store unescaped path
		unescaped := u.Path

		u.Path = ""
		u.RawPath = ""

		return unescaped, nil
	}

	unescapedNext, err := url.PathUnescape(escapedPathParts[1])
	if nil != err {
		return "", err
	}

	u.Path = unescapedNext

	if unescapedNext != escapedPathParts[1] {
		// path requires escaping
		u.RawPath = escapedPathParts[1]
	} else {
		u.RawPath = ""
	}

	return url.PathUnescape(escapedPathParts[0])
}
