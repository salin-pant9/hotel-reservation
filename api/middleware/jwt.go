package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/salin-pant9/hotel-reservation/db"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {

		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			fmt.Println("Token not present in header")
			return fmt.Errorf("unauthorized")
		}
		claims, err := validateToken(token[0])
		if err != nil {
			return err
		}
		expires := int64(claims["expires"].(float64))
		if time.Now().Unix() > expires {
			return fmt.Errorf("token Expired")
		}
		userID := claims["id"].(string)
		user, err := userStore.GetUserById(c.Context(), userID)
		if err != nil {
			return fmt.Errorf("Unauthorized")
		}
		c.Context().SetUserValue("user", user)
		// fmt.Println("token: ", token)
		return c.Next()
	}
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Invalid signing Method", token.Header["alg"])
			return nil, fmt.Errorf("Unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("Failed to parse JWT-Token", err)
		return nil, fmt.Errorf("unauthorized")
	}
	if !token.Valid {
		fmt.Println("Token is not valid")
		return nil, fmt.Errorf("Unauthorized")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Unauthorized")

	}

	return claims, nil
}
