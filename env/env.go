package env

import (
	"log"
	"os"
)

const (
	USER_ID  string = "BASIC_AUTH_USER_ID"
	PASSWORD string = "BASIC_AUTH_PASSWORD"
)

type Env struct {
	UserID   string
	Password string
}

func LoadEnv() *Env {
	userId := os.Getenv(USER_ID)
	if userId == "" {
		log.Fatalf("env error: %s cannot be empty", USER_ID)
	}
	pass := os.Getenv(PASSWORD)
	if pass == "" {
		log.Fatalf("env error: %s cannot be empty", PASSWORD)
	}

	return &Env{
		UserID:   userId,
		Password: pass,
	}
}
