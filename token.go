package main

import (
	"time"
	"errors"
	"math/rand"
)

var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenNotFound = errors.New("token not found")
)

type TokenManager struct {
	tokens []Token
}

type Token struct {
	User    string
	Token   string
	Expire  int64
}

func (t *TokenManager) Get(token string) (Token, error) {
	now := time.Now().Unix()
	for _, v := range t.tokens {
		if v.Token == token {
			if v.Expire < now {
				return Token{}, ErrTokenExpired
			}

			return v, nil
		}
	}
	return Token{}, ErrTokenNotFound
}

func (t *TokenManager) New(user string) Token {
	token := Token {
		user, RandStringRunes(32), time.Now().Unix() + int64(time.Duration(10) * time.Minute),
	}
	t.tokens = append(t.tokens, token)

	return token
}

func (t *Token) Renew() {
	t.Expire = time.Now().Unix() + int64(time.Duration(10) * time.Minute)
}

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}