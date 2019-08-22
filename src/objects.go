package main

import "go.mongodb.org/mongo-driver/bson/primitive"

//PersonalData describe
type PersonalData struct {
	DocumentID  *primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Name        string              `json:"name" bson:"name"`
	LastName    string              `json:"lastName" bson:"lastName"`
	Phone       string              `json:"phone,omitempty" bson:"phone,omitempty"`
	Email       string              `json:"email,omitempty" bson:"email,omitempty"`
	YearOfBirth int                 `json:"yaerOfBirth,omitempty" bson:"yaerOfBirth,omitempty"`
}
