package auth

import (
	"skillKonnect/app/db"
	"skillKonnect/app/models"

	"golang.org/x/crypto/bcrypt"
)

// Event name constants
const (
	UserSignupEvent         = "auth.signup"
	ResendVerificationEvent = "auth.resend.verification"
)

// UserWithVerificationToken is a struct that will be sent over the
// auth.signup event. It holds the User struct and the Verification token string.
type UserWithVerificationToken struct {
	User  models.User
	Token string
}

// type Auth struct {
// 	UserID   uint
// 	Email    string
// 	LoggedIn bool
// }

// func (auth Auth) Check() bool {
// 	return auth.LoggedIn
// }

func createUserFromFormValues(values SignupFormValues) (models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(values.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	user := models.User{
		Email:        values.Email,
		FirstName:    values.FirstName,
		LastName:     values.LastName,
		PasswordHash: string(hash),
	}
	result := db.Get().Create(&user)
	return user, result.Error
}
