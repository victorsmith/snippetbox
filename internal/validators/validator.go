package validators

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

// If the # of errors == 0 => return true
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
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

// PermittedInt() returns true if a value is in a list of permitted integers.
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
