package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
  password := RandomString(10)
  hashed,err:= HashPassword(password)
  require.NoError(t, err)
  require.NotEmpty(t, hashed)
  
  err=CheckPassword(password, hashed)
  require.NoError(t, err)
  
  wrongPass := RandomString(6)
  err = CheckPassword(wrongPass, hashed)
  require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

  secondHash, err := HashPassword(password)
  require.NoError(t, err)
  require.NotEmpty(t, secondHash)
  require.NotEqual(t, hashed, secondHash)
}
