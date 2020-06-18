package jwt

import (
	"testing"
	"time"
)

var (
	secret = []byte("123456")
)

type Payload struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
}

func TestGenerateToken(t *testing.T) {
	payload := &Payload{
		UserID:   101,
		UserName: "hello",
	}
	token, err := GenerateToken(secret, payload, time.Hour)
	if err != nil {
		t.Error(err)
		t.Failed()
		return
	}
	t.Log("token:", token)
	p, err := ParseToken(token, secret)
	if err != nil {
		t.Error(err)
		t.Failed()
		return
	}
	t.Log(p)
}
