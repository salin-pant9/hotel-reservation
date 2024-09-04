package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/salin-pant9/hotel-reservation/api"
	"github.com/salin-pant9/hotel-reservation/api/middleware"
	"github.com/salin-pant9/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// const dbUri = "mongodb+srv://scorpiopant:6m2NT2ThmbFVmm79@cluster0.mkacidw.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

// const dbname = "hotel_reservation"

// const userColl = "users"

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of API server")
	flag.Parse()
	var (

		// Initialize server
		app = fiber.New(config)
		// handler Initialization
		userHandler    = api.NewUserHandler(db.NewMongoUserStore(client, db.DBNAME))
		hotelStore     = db.NewMongoHotelStore(client)
		roomStore      = db.NewMongoRoomStore(client, hotelStore)
		hotelHandler   = api.NewHotelHandler(hotelStore, roomStore)
		userStore      = db.NewMongoUserStore(client, db.DBNAME)
		bookingStore   = db.NewMongoBookingStore(client)
		roomHandler    = api.NewRoomHandler(roomStore, bookingStore)
		authHandler    = api.NewAuthHandler(userStore)
		bookingHandler = api.NewBookingHandler(roomStore, bookingStore)
		// To group apis
		appAuth = app.Group("/api")
		appv1   = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
		admin   = appv1.Group("/admin", middleware.AdminAuth)
	)
	// auth router
	appAuth.Post("/login", authHandler.HandleLogin)

	//user routers
	appv1.Post("/user", userHandler.HandlePostUser)
	appv1.Get("/user", userHandler.HandleGetUsers)
	appv1.Get("/user/:id", userHandler.HandleGetUser)
	appv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	appv1.Put("/user/:id", userHandler.HandlePutUser)

	// hotel routers
	appv1.Get("/hotel", hotelHandler.HandleGetHotel)
	appv1.Get("/hotel/:id", hotelHandler.HandleGetHotelByID)
	appv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// room routers
	appv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	// booking routers
	appv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	// admin routers
	admin.Get("/booking", bookingHandler.HandleGetBookings)
	// server listens in given port
	app.Listen(*listenAddr)
}

// func handleUser(c *fiber.Ctx) error {
// 	return c.JSON(map[string]string{"user": "JOhn Doe"})
// }

// Database URL -> mongodb+srv://scorpiopant:6m2NT2ThmbFVmm79@cluster0.mkacidw.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0
