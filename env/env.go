package env

import (
	"fmt"
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

func GetEnv() (*Env, error) {
	userId := os.Getenv(USER_ID)
	if userId == "" {
		return nil, &EnvError{EnvName: USER_ID}
	}
	pass := os.Getenv(PASSWORD)
	if pass == "" {
		return nil, &EnvError{EnvName: PASSWORD}
	}

	return &Env{
		UserID:   userId,
		Password: pass,
	}, nil
}

type EnvError struct {
	EnvName string
}

func (err *EnvError) Error() string {
	return fmt.Sprintf("env error: %s cannot be empty", err.EnvName)
}
