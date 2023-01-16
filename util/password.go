package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
  hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err!= nil {
    return "", fmt.Errorf("Could not encypt password %s", err.Error())
  }
  return string(hashed), nil
}

func CheckPassword(password, hashedPassword string) error {
  return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
