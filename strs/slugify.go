package strs

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mozillazg/go-unidecode"
)

var (
	// Separator separator between words
	Separator = "-"

	// SeparatorForRe for regexp
	SeparatorForRe = regexp.QuoteMeta(Separator)

	// ReInValidChar match invalid slug string
	ReInValidChar = regexp.MustCompile(fmt.Sprintf("[^%sa-zA-Z0-9]", SeparatorForRe))

	// ReDupSeparatorChar match duplicate separator string
	ReDupSeparatorChar = regexp.MustCompile(fmt.Sprintf("%s{2,}", SeparatorForRe))
)

// Slugify implements make a pretty slug from the given text.
// e.g. Slugify("kožušček hello world") => "kozuscek-hello-world"
func Slugify(s string) string {
	s = unidecode.Unidecode(s)
	s = ReInValidChar.ReplaceAllString(s, Separator)
	s = ReDupSeparatorChar.ReplaceAllString(s, Separator)
	s = strings.Trim(s, Separator)
	s = strings.ToLower(s)
	return s
}
