package models

import (
	"github.com/anthdm/superkit/kit"
)

// AuthPayload represents an user that might be authenticated.
type AuthPayload struct {
	Authenticated bool
	User          *User
	Token         string
}

// Check should return true if the user is authenticated.
// See handlers/auth.go.
func (user AuthPayload) Check() bool {
	return user.User != nil && user.User.ID > 0 && user.Authenticated
}

func (a AuthPayload) GetUser() *User {
	return a.User
}

func (a AuthPayload) GetToken() string {
	return a.Token
}

func (a AuthPayload) IsAdmin() bool {
	return a.User != nil && a.User.Role == "admin"
}

type ExtendedAuth interface {
	kit.Auth // this contains only Check()
	GetUser() *User
	GetToken() string
	IsAdmin() bool
}
