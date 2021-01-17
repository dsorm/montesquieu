package users

import (
	"math/rand"
	"testing"
)
import "github.com/dchest/uniuri"

func TestPasswordHashing(t *testing.T) {
	stdlen := 0
	hash := ""
	passwordValid := false
	var err error
	for i := 0; i < 20; i++ {
		// generate a random password
		stdlen = int(rand.Intn(10) + 5)
		password := uniuri.NewLen(stdlen)

		// test the hashing function
		hash, err = HashPassword(password)
		if err != nil {
			t.Error("HashPassword() returned non nil err value:", err.Error())
		} else if len(hash) == 0 {
			t.Error("HashPassword() returned string hash with len() == 0")
		}

		// test verifying function with the generated hash
		err = nil
		passwordValid, err = VerifyPassword(hash, password)
		if !passwordValid {
			t.Error("VerifyPassword() evaluated valid password as invalid.")
		}
		if err != nil {
			t.Error("VerifyPassword() returned non nil err value:", err.Error())
		}

		// test verifying function with an invalid password
		err = nil
		passwordValid, err = VerifyPassword(hash, "D*23snHdasdsadasdasde22$%")
		if passwordValid {
			t.Error("VerifyPassword() evaluated invalid password as valid.")
		}
		// the function returns an error when the password doesn't work, and that's ok
		if err != nil && err.Error() != "Password did not match" {
			t.Error("VerifyPassword() returned unexpected non nil err value:", err.Error())
		}

	}
	// some edge cases
	hash, err = HashPassword("")
	if err == nil {
		t.Errorf("HashPassword() should return an error with input string of len() == 0")
	}

	passwordValid, err = VerifyPassword("", "dsdsd")
	if passwordValid {
		t.Errorf("VerifyPassword() should always return false for empty hashes")
	}
}
