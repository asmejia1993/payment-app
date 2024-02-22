package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	t.Run("check ok password", func(t *testing.T) {
		password := RandomString(10)
		hashPass1, err := HashPassword(password)
		require.NoError(t, err)
		require.NotEmpty(t, hashPass1)

		err = CheckPassword(hashPass1, password)
		require.NoError(t, err)
	})

	t.Run("check wrong password", func(t *testing.T) {
		pass := RandomString(9)
		wrongPass := RandomString(10)
		hashWrong, _ := HashPassword(pass)
		err := CheckPassword(hashWrong, wrongPass)
		require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
	})
}
