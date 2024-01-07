package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/amitamrutiya/09-go-mongodb-api/models"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UserController represents the controller for operating on the User resource
type UserController struct {
	session *mgo.Session
}

// NewUserController provides a reference to a UserController with provided mongo session
func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

// GetUser retrieves an individual user resource
func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId hex representation, otherwise return status not found
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Stub user
	u := models.User{}

	// Fetch user
	if err := uc.session.DB("09-go-mongodb-api").C("users").FindId(oid).One(&u); err != nil {
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	// Marshal provided interface into JSON structure
	uj, err := json.Marshal(u)

	// Write error
	if err != nil {
		panic(err)
	}

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	w.Write(uj)
}

// CreateUser creates a new user resource
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Stub an user to be populated from the body
	u := models.User{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	// Add an Id
	u.ID = bson.NewObjectId()

	// Write the user to mongo
	uc.session.DB("09-go-mongodb-api").C("users").Insert(u)

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	w.Write(uj)
}

// DeleteUser removes an existing user resource
func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId hex representation, otherwise return status not found
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Remove user
	if err := uc.session.DB("09-go-mongodb-api").C("users").RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	// Write status
	w.WriteHeader(http.StatusOK) // 200
}

// UpdateUser updates an existing user resource
func (uc UserController) UpdateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId hex representation, otherwise return status not found
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	// Stub user
	u := models.User{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Update user
	if err := uc.session.DB("09-go-mongodb-api").C("users").UpdateId(oid, u); err != nil {
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	// Write status
	w.WriteHeader(http.StatusOK) // 200
}
