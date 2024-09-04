package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/salin-pant9/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	hotelStore db.HotelStore
	roomStore  db.RoomStore
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore:  rs,
	}
}
func (hh *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"hotelID": oid}
	rooms, err := hh.roomStore.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (hh *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	hotels, err := hh.hotelStore.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}
func (hh *HotelHandler) HandleGetHotelByID(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	hotel, err := hh.hotelStore.GetHotelByID(c.Context(), oid)
	if err != nil {
		return err
	}
	return c.JSON(hotel)
}
