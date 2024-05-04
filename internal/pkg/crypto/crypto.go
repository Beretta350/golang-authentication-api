package crypto

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) ([]byte, error) {
	encryptedData, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password: ", err)
		return nil, err
	}
	return encryptedData, nil
}

func CheckPassword(password string, hash []byte) error {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		log.Println("Password does not match hash: ", err)
		return err
	}
	return nil
}
