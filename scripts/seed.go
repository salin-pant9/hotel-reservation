package main

import (
	"context"
	"log"

	"github.com/salin-pant9/hotel-reservation/db"
	"github.com/salin-pant9/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx        = context.Background()
	client     *mongo.Client
	hotelStore db.HotelStore
	roomStore  db.RoomStore
	userStore  db.UserStore
)

func seedUser(isAdmin bool, fname, lname, email string) {
	user, err := types.NewUserFromParams(types.CreateUserParam{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  "123456",
	})
	if err != nil {
		log.Fatal(err)

	}
	user.IsAdmin = isAdmin
	_, err = userStore.PostUsers(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
}

func seedHotel(name, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	rooms := []types.Room{
		{
			Type:      types.SingleRoomType,
			BasePrice: 99.99,
		},
		{
			Type:      types.DeluxeRoomType,
			BasePrice: 1200,
		},
		{
			Type:      types.SeaSideRoomType,
			BasePrice: 1000,
		},
	}
	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {

		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	seedHotel("Shangrila", "Kathmandu", 3)
	seedUser(true, "james", "cram", "james@gmail.com")

}

func init() {

	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client, db.DBNAME)
}
