package cache

import "strings"

var AllowedTypes = []string{
	"image/gif",
	"text/html",
}

var BlockedTypes = []string{
	"video/mpeg",
}

func MayBeCached(ctype string) bool {
	for _, t := range BlockedTypes {
		if strings.Contains(strings.ToLower(ctype), strings.ToLower(t)) {
			return false
		}
	}

	for _, t := range AllowedTypes {
		if strings.Contains(strings.ToLower(ctype), strings.ToLower(t)) {
			return true
		}

	}

	return false
}
