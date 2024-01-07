package models

import "gopkg.in/mgo.v2/bson"

type User struct {
	ID       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Username string        `json:"username,omitempty" bson:"username,omitempty"`
	Gender   string        `json:"gender" bson:"gender"`
	Age      int           `json:"age" bson:"age"`
}
