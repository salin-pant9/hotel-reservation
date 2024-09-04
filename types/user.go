package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
)

type CreateUserParam struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (cup CreateUserParam) Validate() map[string]string {
	errors := map[string]string{}
	if len(cup.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("firstName length should be of %d characters", minFirstNameLen)
	}
	if len(cup.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("lastName length should be of %d characters", minLastNameLen)
	}
	if len(cup.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password length should be of %d characters", minPasswordLen)
	}
	if !ValidateEmail(cup.Email) {
		errors["email"] = fmt.Sprintf("email is not valid")
	}
	return errors
}

func ValidateEmail(e string) bool {
	emailRegx := regexp.MustCompile(`^([a-zA-Z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4})$`)
	return emailRegx.MatchString(e)
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"first_name" json:"firstName"`
	LastName          string             `bson:"last_name" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"EncryptedPassword" json:"-"`
	IsAdmin           bool               `bson:"IsAdmin" json:"IsAdmin"`
}

func NewUserFromParams(params CreateUserParam) (*User, error) {
	enpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(enpw),
	}, nil
}

func IsValidPassword(encryptedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password)) == nil
}
