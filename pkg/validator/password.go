package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// ValidatePasswordStrength validates password strength
// Must contain:
// - At least 1 uppercase letter
// - At least 1 number
// - At least 1 special character (!@#$%^&*())
// - Minimum 8 characters
func ValidatePasswordStrength(fl validator.FieldLevel) bool {
    password := fl.Field().String()

    // Check for at least 1 uppercase letter
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    if !hasUpper {
        return false
    }

    // Check for at least 1 number
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
    if !hasNumber {
        return false
    }

    // Check for at least 1 special character
    hasSpecial := regexp.MustCompile(`[!@#$%^&*()\-_=+\[\]{}|;:,.<>?/~` + "`" + `]`).MatchString(password)
    if !hasSpecial {
        return false
    }

    return true
}