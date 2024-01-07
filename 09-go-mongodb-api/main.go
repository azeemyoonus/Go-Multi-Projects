package main

import (
	"log"
	"net/http"

	"github.com/amitamrutiya/09-go-mongodb-api/controllers"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

func main() {

	r := httprouter.New()

	// Get a UserController instance
	uc := controllers.NewUserController(getSession())
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)
	r.PUT("/user/:id", uc.UpdateUser)

	// Fire up the server
	http.ListenAndServe("localhost:3000", r)
}

func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb+srv://akamrutiya22102002:mongopassword@cluster0.99vsucr.mongodb.net/")

	// Check if connection error, is mongo running?
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		panic(err)
	}

	return s
}
