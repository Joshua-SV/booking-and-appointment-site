package auth

import "golang.org/x/crypto/bcrypt"

func HashedPassword(password string) (string, error) {
	// get the hash string from bcrypt
	passHashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(passHashed), nil
}
