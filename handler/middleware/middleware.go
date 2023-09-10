package middleware

import "github.com/TechBowl-japan/go-stations/env"

type Middleware struct {
	env *env.Env
}

func NewMiddleware(env *env.Env) *Middleware {
	return &Middleware{env}
}
