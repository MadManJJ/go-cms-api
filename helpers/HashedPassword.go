package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

func GetHashed(password string) []byte {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)	

	return hashed
}