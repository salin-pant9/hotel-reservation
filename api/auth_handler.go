package api

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/salin-pant9/hotel-reservation/db"
	"github.com/salin-pant9/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

func (uh *AuthHandler) HandleLogin(c *fiber.Ctx) error {
	var params AuthParams

	if err := c.BodyParser(&params); err != nil {
		return err
	}
	user, err := uh.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("Invalid Credentials")
		}
		return err
	}
	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return fmt.Errorf("Invalid Credentials")
	}
	resp := LoginResponse{
		User:  user,
		Token: CreateTokenFromUser(user),
	}

	return c.JSON(resp)
}
func CreateTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret key")
	}
	return tokenStr
}
