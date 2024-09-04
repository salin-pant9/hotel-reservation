package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/salin-pant9/hotel-reservation/db"
	"github.com/salin-pant9/hotel-reservation/types"
)

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParam{
		FirstName: "James",
		LastName:  "Foo",
		Email:     "foo@gmail.com",
		Password:  "123456",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = userStore.PostUsers(context.TODO(), user)
	if err != nil {
		t.Fatal(err)

	}
	return user
}

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertedUser := insertTestUser(t, tdb.UserStore)
	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleLogin)

	params := AuthParams{
		Email:    "foo@gmail.com",
		Password: "123456",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)

	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 but got %d", resp.StatusCode)
	}
	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatal(err)
	}
	if loginResp.Token == "" {
		t.Fatalf("Expected token to be present")
	}
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, loginResp.User) {
		t.Fatal("Inserted User should be the returned User")
	}
	fmt.Println(resp)
}
