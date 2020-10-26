package strs

import "regexp"

var (
	emailMatcher = regexp.MustCompile(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`)
	phoneMatcher = regexp.MustCompile(`^1[3-9]\d{9}$`)
)

func IsEmail(s string) bool {
	return emailMatcher.MatchString(s)
}

func IsMobilePhone(s string) bool {
	return phoneMatcher.MatchString(s)
}
