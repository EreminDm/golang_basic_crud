package controller

import (
	"context"

	"github.com/EreminDm/golang_basic_crud/database"
	"github.com/pkg/errors"
)

// PersonalData description.
type PersonalData struct {
	DocumentID  string
	Name        string
	LastName    string
	Phone       string
	Email       string
	YearOfBirth int
}

// UsersPersonalData description.
type UsersPersonalData interface {
	One(ctx context.Context, value string) (PersonalData, error)
	All(ctx context.Context) ([]PersonalData, error)
	Remove(ctx context.Context, id string) (int64, error)
	Update(ctx context.Context, p *PersonalData) (int64, error)
	Insert(ctx context.Context, document *PersonalData) (interface{}, error)
}

// Insert adds data to collection.
func (p PersonalData) Insert(ctx context.Context, document *PersonalData) (interface{}, error) {
	var u database.User
	var doc = database.PersonalData{DocumentID: p.DocumentID, Email: p.Email, LastName: p.LastName, Name: p.Name, Phone: p.Phone, YearOfBirth: p.YearOfBirth}
	i, err := u.Insert(ctx, &doc)
	if err != nil {
		return nil, errors.Wrap(err, "could not insert personal data")
	}
	return i, nil
}

// One returns personal data from collection.
func (p PersonalData) One(ctx context.Context, value string) (*PersonalData, error) {
	var u database.User
	user, err := u.One(ctx, value)
	if err != nil {
		return nil, errors.Wrap(err, "could not select one personal data")
	}

	p.DocumentID = user.DocumentID
	p.Email = user.Email
	p.LastName = user.LastName
	p.Name = user.Name
	p.Phone = user.Phone
	p.YearOfBirth = user.YearOfBirth

	return &p, nil
}

// All returns an array of personal information.
func (p PersonalData) All(ctx context.Context) ([]PersonalData, error) {
	var (
		u  database.User
		pr []PersonalData
	)
	users, err := u.All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "could not select all personal data")
	}
	for _, usr := range users {
		p.DocumentID = usr.DocumentID
		p.Email = usr.Email
		p.LastName = usr.LastName
		p.Name = usr.Name
		p.Phone = usr.Phone
		p.YearOfBirth = usr.YearOfBirth
		pr = append(pr, p)
	}
	return pr, nil
}

// Update change information in collection.
func (p PersonalData) Update(ctx context.Context, pd *PersonalData) (int64, error) {
	var u database.User
	var doc = database.PersonalData{DocumentID: pd.DocumentID, Email: pd.Email, LastName: pd.LastName, Name: pd.Name, Phone: pd.Phone, YearOfBirth: pd.YearOfBirth}
	return u.Update(ctx, &doc)
}

// Remove function deletes information from collection.
func (p *PersonalData) Remove(ctx context.Context, id string) (int64, error) {
	var u database.User
	return u.Remove(ctx, id)
}
