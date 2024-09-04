package db

import (
	"context"
	"fmt"

	"github.com/salin-pant9/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "users"

type Dropper interface {
	Drop(context.Context) error
}

type UserStore interface {
	Dropper
	GetUserByEmail(context.Context, string) (*types.User, error)
	GetUserById(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	PostUsers(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, bson.M, bson.M) error
}
type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client, dbname string) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(dbname).Collection(userColl),
	}
}
func (mus *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("------Dropping collection")
	return mus.coll.Drop(ctx)
}
func (mus *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user types.User
	if err := mus.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
func (mus *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = mus.coll.DeleteOne(ctx, bson.M{"_id": oId})
	if err != nil {
		return err
	}
	return nil
}
func (mus *MongoUserStore) UpdateUser(ctx context.Context, filter, values bson.M) error {
	update := bson.M{"$set": values}
	fmt.Println(values)
	_, err := mus.coll.UpdateOne(ctx, update, filter)
	if err != nil {
		return err
	}
	return nil
}
func (mus *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	var users []*types.User
	cur, err := mus.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}
func (mus *MongoUserStore) PostUsers(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := mus.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}
func (mus *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User
	if err := mus.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
