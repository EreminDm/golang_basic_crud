package crud

import (
	"context"
	"github.com/EreminDm/golang_basic_crud/database"
)

// PersonalData description.
type PersonalData struct {
	DocumentID  [12]byte `json:"id,omitempty" bson:"_id"` // as *primitive.ObjectID
	Name        string   `json:"name" bson:"name"`
	LastName    string   `json:"lastName" bson:"lastName"`
	Phone       string   `json:"phone,omitempty" bson:"phone,omitempty"`
	Email       string   `json:"email,omitempty" bson:"email,omitempty"`
	YearOfBirth int      `json:"yaerOfBirth,omitempty" bson:"yaerOfBirth,omitempty"`
}

// UsersPersonalData abc
type UsersPersonalData interface {
	SelectOne(ctx context.Context, key, value string) (result PersonalData, err error)
	SelectAll(ctx context.Context) (results []PersonalData, err error)
	Remove(ctx context.Context, id string) (int64, error)
	Update(ctx context.Context, p *PersonalData) (int64, error)
	Insert(ctx context.Context, document *PersonalData) (interface{}, error)
}

func (p  PersonalData) SelectOne(ctx context.Context, key, value string) (result PersonalData, err error){
	return database.SelectOne(ctx,key,value)
}
