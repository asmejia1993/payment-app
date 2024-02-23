package token

import (
	"testing"
	"time"

	"github.com/asmejia1993/payment-app/db/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateToken_NewPayload(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	email := "test@example.com"
	duration := time.Hour

	token, payload, err := maker.CreateToken(email, duration)

	assert.NoError(t, err)
	assert.NotNil(t, payload)
	assert.NotEmpty(t, payload.ID)
	assert.Equal(t, email, payload.Email)
	assert.WithinDuration(t, time.Now(), payload.IssuedAt, time.Second)
	assert.WithinDuration(t, time.Now().Add(duration), payload.ExpiredAt, time.Second)
	assert.NotEmpty(t, token)
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

func TestCreateToken_ExpiredTimeInFuture(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	email := "test@example.com"
	duration := time.Hour

	token, payload, err := maker.CreateToken(email, duration)

	assert.NoError(t, err)
	assert.NotNil(t, payload)
	assert.NotEmpty(t, payload.ID)
	assert.Equal(t, email, payload.Email)
	assert.WithinDuration(t, time.Now(), payload.IssuedAt, time.Second)
	assert.WithinDuration(t, time.Now().Add(duration), payload.ExpiredAt, time.Second)
	assert.NotEmpty(t, token)
}
