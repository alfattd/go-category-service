package validator

import (
	"strings"
	"unicode"
)

type ErrorsValidator struct {
	Messages []string
}

func (v *ErrorsValidator) Error() string {
	return strings.Join(v.Messages, "; ")
}

func (v *ErrorsValidator) Add(msg string) {
	v.Messages = append(v.Messages, msg)
}

func (v *ErrorsValidator) HasErrors() bool {
	return len(v.Messages) > 0
}

const (
	MaxCategoryNameLength = 20
)

var forbiddenRunes = map[rune]bool{
	'<': true, '>': true,
	'"': true, '\'': true,
	';': true, '&': true,
	'\\': true, '/': true,
	'{': true, '}': true,
	'(': true, ')': true,
	'[': true, ']': true,
}

func CategoryNameValidator(name string) *ErrorsValidator {
	errs := &ErrorsValidator{}

	if name == "" {
		errs.Add("name is required")
		return errs
	}

	if len([]rune(name)) > MaxCategoryNameLength {
		errs.Add("name must not exceed 20 characters")
	}

	if hasForbiddenRunes(name) {
		errs.Add("name contains invalid characters (< > \" ' ; & \\ / { } ( ) [ ] are not allowed)")
	}

	if hasOnlyWhitespace(name) {
		errs.Add("name must not be blank")
	}

	if errs.HasErrors() {
		return errs
	}

	return nil
}

func hasForbiddenRunes(s string) bool {
	for _, r := range s {
		if forbiddenRunes[r] {
			return true
		}
	}
	return false
}

func hasOnlyWhitespace(s string) bool {
	return strings.TrimFunc(s, unicode.IsSpace) == ""
}
