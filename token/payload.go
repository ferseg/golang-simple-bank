package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
  ErrExpiredToken = errors.New("Token has expired")
  ErrAlgNotAllowed = errors.New("The algorithm is not allowed")
  ErrInvalidToken = errors.New("The token is invalid")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
  tokenId, err := uuid.NewRandom()
  if err != nil {
    return nil, err
  }

  now := time.Now()

  payload := &Payload{
    ID: tokenId,
    Username: username,
    IssuedAt: now,
    ExpiredAt: now.Add(duration),
  }
  return payload, nil
}

func (payload *Payload) Valid() error {
  if time.Now().After(payload.ExpiredAt) {
    return ErrExpiredToken
  }
  return nil
}
