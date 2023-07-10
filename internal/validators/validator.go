package validators

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Returns a pointer to a compiled regexp.Regexp type
var EmailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// If the # of errors (field & non field) == 0 => return true
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// AddFieldError() adds an error message to the FieldErrors map (so long as no
// entry already exists for the given key).
func (v *Validator) AddFieldError(key, error string) {
	// If no errors exist, init a maps
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	// If no error exists at "key" => add one
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = error
	}
}

func (v *Validator) AddNonFieldError(error string) {
	v.NonFieldErrors = append(v.NonFieldErrors, error)
}

// CheckField() adds an error message to the FieldErrors map only if a
// validation check is not 'ok'.
func (v *Validator) CheckField(ok bool, key, error string) {
	if !ok {
		v.AddFieldError(key, error)
	}
}

// NotBlank() returns true if a value is not an empty string
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars() returns true if a value contains no more than "max" characters.
func MaxChars(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}

// PermittedValue() returns true if a value is in a list of permitted comparable types (generic).
func PermittedValue[T comparable] (value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

// Returns true if the input string contains at least "n" letters
func MinChars(input string, n int) bool {
	return utf8.RuneCountInString(input) >= n
}

// Returns true if input matches the provided regexp (compiled)
func Matches(input string, rx *regexp.Regexp) bool {
	return rx.MatchString(input)
}
