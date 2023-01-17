package token

import (
	"testing"
	"time"

	"github.com/ferseg/golang-simple-bank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {

  maker, err := NewPasetoMaker(util.RandomString(32))
  require.NoError(t, err)
  
  username := util.RandomOwner()
  duration := time.Minute

  issuedAt := time.Now()
  expiredAt := issuedAt.Add(duration)

  token, err := maker.CreateToken(username, duration)
  require.NoError(t, err)
  require.NotEmpty(t, token)

  paylaod, err := maker.VerifyToken(token)
  require.NoError(t, err)
  require.NotEmpty(t, paylaod)

  require.NotZero(t, paylaod.ID)
  require.Equal(t, username, paylaod.Username)
  require.WithinDuration(t, issuedAt, paylaod.IssuedAt, time.Second)
  require.WithinDuration(t, expiredAt, paylaod.ExpiredAt, time.Second)
  
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

  payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
