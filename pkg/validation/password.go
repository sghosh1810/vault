package validation

import (
	"errors"
	"unicode"
)

var (
	ErrPasswordTooShort      = errors.New("password must be at least 8 characters long")
	ErrPasswordTooLong       = errors.New("password must be at most 64 characters long")
	ErrPasswordNoLower       = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoUpper       = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoDigit       = errors.New("password must contain at least one digit")
	ErrPasswordNoSpecial     = errors.New("password must contain at least one special character")
	ErrPasswordHasWhitespace = errors.New("password must not contain spaces")
)

// ValidatePassword enforces an industry-standard strong password policy.
func ValidatePassword(password string) error {
	const minLen = 8
	const maxLen = 64

	if len(password) < minLen {
		return ErrPasswordTooShort
	}
	if len(password) > maxLen {
		return ErrPasswordTooLong
	}

	var hasLower, hasUpper, hasDigit, hasSpecial bool

	for _, c := range password {
		switch {
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		case unicode.IsSpace(c):
			return ErrPasswordHasWhitespace
		}
	}

	if !hasLower {
		return ErrPasswordNoLower
	}
	if !hasUpper {
		return ErrPasswordNoUpper
	}
	if !hasDigit {
		return ErrPasswordNoDigit
	}
	if !hasSpecial {
		return ErrPasswordNoSpecial
	}

	return nil
}
