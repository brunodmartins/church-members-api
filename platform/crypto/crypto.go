package crypto

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) []byte {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		logrus.Panicf("Error encrypting password:%v", err)
	}
	return encryptedPassword
}

func IsSamePassword(encryptedPassword []byte, rawPassword string) error {
	return bcrypt.CompareHashAndPassword(encryptedPassword, []byte(rawPassword))
}

