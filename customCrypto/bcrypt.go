// Pacakge customcrypto stores all the constants used all over the service

package customcrypto

import "golang.org/x/crypto/bcrypt"

type Bcrypt struct {
	Cost int
}

func (bc *Bcrypt) IsSame(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func (bc *Bcrypt) GetHash(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bc.Cost)
}

func (bc *Bcrypt) NoMatch(err error) bool {
	return err == bcrypt.ErrMismatchedHashAndPassword
}
