package token

import (
	"reflect"
	"testing"
	"time"

	"github.com/asmejia1993/payment-app/db/util"
	"github.com/o1egl/paseto"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	email := "test@example.com"
	duration := time.Hour

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(email, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, email, payload.Email)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestPasetoMaker_CreateToken(t *testing.T) {
	type fields struct {
		paseto       *paseto.V2
		symmetricKey []byte
	}
	type args struct {
		email    string
		duration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   *Payload
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maker := &PasetoMaker{
				paseto:       tt.fields.paseto,
				symmetricKey: tt.fields.symmetricKey,
			}
			got, got1, err := maker.CreateToken(tt.args.email, tt.args.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("PasetoMaker.CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PasetoMaker.CreateToken() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PasetoMaker.CreateToken() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
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
