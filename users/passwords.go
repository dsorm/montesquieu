package users

import (
	"github.com/raja/argon2pw"
)

// HashPassword hashes password from user's input, so it can be safely put
// into a database
func HashPassword(input string) (string, error) {
	hashedPassword, err := argon2pw.GenerateSaltedHash(input)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

// VerifyPassword verifies the password that has user provided ('input') against
// the hash
// Returns true if the password does match, and false if it doesn't or if error
// has occured during the check
func VerifyPassword(hash string, input string) (bool, error) {

	// Don't allow login for users without any password hash
	if hash == "" {
		return false, nil
	}

	passwordValid, err := argon2pw.CompareHashWithPassword(hash, input)
	if err != nil {
		return false, err
	}

	return passwordValid, nil
}
