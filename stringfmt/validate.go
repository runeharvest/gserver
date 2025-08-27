package stringfmt

import (
	"fmt"
	"strings"
)

func PasswordValidate(password string) error {

	return nil
}

func UsernameValidate(username string) error {
	profaneWords := []string{"damn", "hell", "shit", "fuck"} // Add more as needed
	for _, word := range profaneWords {
		if len(username) < len(word) {
			continue
		}
		if !containsInsensitive(username, word) {
			continue
		}
		return fmt.Errorf("contains inappropriate language")
	}

	return nil
}
func containsInsensitive(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
