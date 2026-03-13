package slug

import "strings"

// Slugify normalizes a string into a lowercase, hyphen-separated slug.
func Slugify(s string) string {
	s = strings.ToLower(s)

	var b strings.Builder
	lastHyphen := false

	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
			lastHyphen = false
		case r >= '0' && r <= '9':
			b.WriteRune(r)
			lastHyphen = false
		case r == '-' || r == ' ':
			if b.Len() == 0 || lastHyphen {
				continue
			}
			b.WriteByte('-')
			lastHyphen = true
		}
	}

	return strings.Trim(b.String(), "-")
}
