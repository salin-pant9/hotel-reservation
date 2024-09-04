package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/salin-pant9/hotel-reservation/db"
	"github.com/salin-pant9/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomsParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

func (brp BookRoomsParams) Validate() error {
	now := time.Now()
	if now.After(brp.FromDate) || now.After(brp.TillDate) {
		return fmt.Errorf("Cannot book hotel")
	}
	return nil
}

type RoomHandler struct {
	RoomStore    db.RoomStore
	BookingStore db.BookingStore
}

func NewRoomHandler(roomStore db.RoomStore, bookingStore db.BookingStore) *RoomHandler {
	return &RoomHandler{
		RoomStore:    roomStore,
		BookingStore: bookingStore,
	}
}

func (rh *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomsParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.Validate(); err != nil {
		return err
	}
	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(map[string]string{"err": "Internal Server Error"})
	}
	ok, err = rh.IsRoomAvailableForBooking(c.Context(), params, roomID)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Booking already made"})
	}
	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}
	inserted, err := rh.BookingStore.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}
	// fmt.Println(booking)

	return c.JSON(inserted)
}

func (rh *RoomHandler) IsRoomAvailableForBooking(ctx context.Context, params BookRoomsParams, roomID primitive.ObjectID) (bool, error) {
	where := bson.M{
		"roomID": roomID,
		"FromDate": bson.M{
			"$gte": params.FromDate,
		},
		"TillDate": bson.M{
			"$lte": params.TillDate,
		},
	}
	bookings, err := rh.BookingStore.GetBookings(ctx, where)
	if err != nil {
		return false, err
	}
	ok := len(bookings) == 0
	return ok, nil

}
