package token

import (
	"testing"
	"time"

	"github.com/ferseg/golang-simple-bank/util"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
  maker, err := NewJWTMaker(util.RandomString(32))
  require.NoError(t, err)
  
  username := util.RandomOwner()
  duration := time.Minute

  issuedAt := time.Now()
  expiredAt := issuedAt.Add(duration)

  token, payload, err := maker.CreateToken(username, duration)
  require.NoError(t, err)
  require.NotEmpty(t, token)
  require.NotEmpty(t, payload)

  paylaod, err := maker.VerifyToken(token)
  require.NoError(t, err)
  require.NotEmpty(t, paylaod)

  require.NotZero(t, paylaod.ID)
  require.Equal(t, username, paylaod.Username)
  require.WithinDuration(t, issuedAt, paylaod.IssuedAt, time.Second)
  require.WithinDuration(t, expiredAt, paylaod.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
  require.NotEmpty(t, payload)

  payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrAlgNotAllowed.Error())
	require.Nil(t, payload)
}

