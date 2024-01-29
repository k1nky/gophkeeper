package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type ID uint64

//go:generate easyjson user.go
//easyjson:json
type User struct {
	ID       ID
	Login    string `json:"login"`
	Password string `json:"password"`
}

type PrivateClaims struct {
	ID    ID
	Login string
}

type contextKey int

const (
	KeyUserClaims contextKey = iota
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) IsValid() error {
	if u.Login == "" || u.Password == "" {
		return ErrCredentialsInvalidFormat
	}
	return nil
}

func NewPrivateClaims(u User) PrivateClaims {
	return PrivateClaims{
		ID:    u.ID,
		Login: u.Login,
	}
}

func NewContextWithClaims(parent context.Context, claims PrivateClaims) context.Context {
	return context.WithValue(parent, KeyUserClaims, claims)
}

func GetEffectiveUser(ctx context.Context) (PrivateClaims, bool) {
	claims, ok := ctx.Value(KeyUserClaims).(PrivateClaims)
	return claims, ok
}
