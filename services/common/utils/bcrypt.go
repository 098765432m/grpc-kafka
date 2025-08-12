package utils

import "golang.org/x/crypto/bcrypt"

func ComparePassword(password []byte, hashedPassword []byte) (bool, error) {

	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		return false, err
	}

	return true, nil
}
