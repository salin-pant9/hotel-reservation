package api

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/salin-pant9/hotel-reservation/db"
	"github.com/salin-pant9/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	RoomStore    db.RoomStore
	BookingStore db.BookingStore
}

func NewBookingHandler(roomStore db.RoomStore, bookingStore db.BookingStore) *BookingHandler {
	return &BookingHandler{
		RoomStore:    roomStore,
		BookingStore: bookingStore,
	}
}

// TODO: this needs to be "admin" authorized
func (bk *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := bk.BookingStore.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

// TODO: this needs to be "user" authorized
func (bk *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := bk.BookingStore.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fmt.Errorf("Not Authorized")
	}
	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(map[string]string{"msg": "Not Authorized"})
	}
	return c.JSON(booking)
}
