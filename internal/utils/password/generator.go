package password

import (
	"crypto/rand"
	"math/big"
)

const (
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	digits           = "0123456789"
	specialChars     = "!@#$%^&*()_+-=[]{}|;':\",./<>?`~"
)

// GenerateTemporaryPassword generates a cryptographically secure temporary password
// The password is 12 characters long and contains uppercase, lowercase, digits, and special characters
func GenerateTemporaryPassword() (string, error) {
	const passwordLength = 12
	charset := uppercaseLetters + lowercaseLetters + digits + specialChars

	password := make([]byte, passwordLength)
	for i := range password {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[randomIndex.Int64()]
	}

	return string(password), nil
}

// GeneratePassword generates a password with specific requirements
func GeneratePassword(length int, requireUpper, requireLower, requireDigits, requireSpecial bool) (string, error) {
	if length < 6 {
		length = 6
	}

	var charset string
	var required []byte

	// Add required character types
	if requireUpper {
		charset += uppercaseLetters
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(uppercaseLetters))))
		if err != nil {
			return "", err
		}
		required = append(required, uppercaseLetters[idx.Int64()])
	}

	if requireLower {
		charset += lowercaseLetters
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(lowercaseLetters))))
		if err != nil {
			return "", err
		}
		required = append(required, lowercaseLetters[idx.Int64()])
	}

	if requireDigits {
		charset += digits
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		required = append(required, digits[idx.Int64()])
	}

	if requireSpecial {
		charset += specialChars
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(specialChars))))
		if err != nil {
			return "", err
		}
		required = append(required, specialChars[idx.Int64()])
	}

	// If no charset specified, use all
	if charset == "" {
		charset = uppercaseLetters + lowercaseLetters + digits + specialChars
	}

	// Generate remaining characters
	remainingLength := length - len(required)
	password := make([]byte, remainingLength)
	for i := range password {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[randomIndex.Int64()]
	}

	// Combine required and random characters
	finalPassword := append(required, password...)

	// Shuffle the password
	for i := len(finalPassword) - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", err
		}
		finalPassword[i], finalPassword[j.Int64()] = finalPassword[j.Int64()], finalPassword[i]
	}

	return string(finalPassword), nil
}
