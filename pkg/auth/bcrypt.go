package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct{}

func (b *Bcrypt) HashAndSalt(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash)
}

func (b *Bcrypt) ComparePassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
