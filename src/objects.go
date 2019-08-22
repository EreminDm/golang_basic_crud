package main

import "gopkg.in/mgo.v2/bson"

//PersonalData describe
type PersonalData struct {
	DocumentID  bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	LastName    string        `json:"lastName" bson:"lastName"`
	Phone       string        `json:"phone,omitempty" bson:"phone,omitempty"`
	Email       string        `json:"email,omitempty" bson:"email,omitempty"`
	YearOfBirth int           `json:"yaerOfBirth,omitempty" bson:"yaerOfBirth,omitempty"`
}
