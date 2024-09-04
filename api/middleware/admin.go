package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/salin-pant9/hotel-reservation/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fmt.Errorf("Not Authorized")
	}
	if !user.IsAdmin {
		return fmt.Errorf("Not Authorized")
	}
	return c.Next()
}
