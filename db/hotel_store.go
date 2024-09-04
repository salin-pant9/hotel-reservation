package db

import (
	"context"

	"github.com/salin-pant9/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotel(context.Context, bson.M, bson.M) error
	GetHotels(context.Context, bson.M) ([]*types.Hotel, error)
	GetHotelByID(context.Context, primitive.ObjectID) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(DBNAME).Collection("hotels"),
	}
}

func (mhs *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := mhs.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID)
	return hotel, nil
}
func (mhs *MongoHotelStore) UpdateHotel(ctx context.Context, filter, update bson.M) error {
	_, err := mhs.coll.UpdateOne(ctx, filter, update)
	return err
}
func (mhs *MongoHotelStore) GetHotels(ctx context.Context, filter bson.M) ([]*types.Hotel, error) {
	res, err := mhs.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel
	if err := res.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}
func (mhs *MongoHotelStore) GetHotelByID(ctx context.Context, id primitive.ObjectID) (*types.Hotel, error) {
	var hotel *types.Hotel
	if err := mhs.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&hotel); err != nil {
		return nil, err
	}
	return hotel, nil
}
