package validation

import (
	"net/mail"
)

// IsValidEmail returns true if the email has a valid RFC5322-like format.
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
