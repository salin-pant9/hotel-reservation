package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/salin-pant9/hotel-reservation/db"
	"github.com/salin-pant9/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbUri = "mongodb+srv://scorpiopant:6m2NT2ThmbFVmm79@cluster0.mkacidw.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
const dbname = "hotel_reservation_test"

type testdb struct {
	UserStore db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
	}
	return &testdb{
		UserStore: db.NewMongoUserStore(client, dbname),
	}
}
func TestPostUser(t *testing.T) {
	tdb := setup(t)
	// fmt.Println(tdb.UserStore)
	defer tdb.teardown(t)
	// t.Fail()
	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)
	params := types.CreateUserParam{
		Email:     "foo@bar.com",
		FirstName: "foobar",
		LastName:  "Carl",
		Password:  "kdfkdjfkdj",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	json.NewDecoder(res.Body).Decode(&user)
	// fmt.Println(user)
	if user.FirstName != params.FirstName {
		t.Errorf("Expected firstname %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("Expected lastname %s but got %s", params.LastName, user.LastName)

	}
	if user.Email != params.Email {
		t.Errorf("Expected email %s but got %s", params.Email, user.Email)
	}

}
