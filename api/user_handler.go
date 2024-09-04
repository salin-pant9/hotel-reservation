package api

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/salin-pant9/hotel-reservation/db"
	"github.com/salin-pant9/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userstore db.UserStore
}

func NewUserHandler(userstore db.UserStore) *UserHandler {
	return &UserHandler{
		userstore: userstore,
	}
}

func (uh *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		values bson.M
		userId = c.Params("id")
	)
	oId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	if err := c.BodyParser(&values); err != nil {
		return err
	}
	filter := bson.M{"_id": oId}
	if err := uh.userstore.UpdateUser(c.Context(), filter, values); err != nil {
		return err
	}
	return c.JSON(map[string]string{"msg": "updated successfully"})
}

func (uh *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParam
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	postedUser, err := uh.userstore.PostUsers(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(postedUser)
}
func (uh *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if err := uh.userstore.DeleteUser(c.Context(), userID); err != nil {
		return err
	}
	return c.JSON(map[string]string{"msg": "deleted successfully"})
}

func (uh *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id  = c.Params("id")
		ctx = context.Background()
	)
	user, err := uh.userstore.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"msg": "not found"})
		}
		return err
	}
	return c.JSON(user)
}

func (uh *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := uh.userstore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}
