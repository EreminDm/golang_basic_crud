package controller

import (
	"context"

	"github.com/EreminDm/golang_basic_crud/entity"
	"github.com/pkg/errors"
)

// Personal describes database implementation.
type Personal struct {
	DB DBProvider
}

// DBProvider describes implementing methods.
type DBProvider interface {
	One(ctx context.Context, value string) (*entity.PersonalData, error)
	All(ctx context.Context) ([]*entity.PersonalData, error)
	Remove(ctx context.Context, id string) (int64, error)
	Update(ctx context.Context, p *entity.PersonalData) (int64, error)
	Insert(ctx context.Context, document *entity.PersonalData) (interface{}, error)
}

// NewPersonal returns new Personal provider.
func NewPersonal(db DBProvider) *Personal {
	return &Personal{
		DB: db,
	}
}

// Insert adds data to collection.
func (p *Personal) Insert(ctx context.Context, document *entity.PersonalData) (interface{}, error) {
	i, err := p.DB.Insert(ctx, document)
	if err != nil {
		return nil, errors.Wrap(err, "could not insert personal data")
	}
	return i, nil
}

// One returns personal data from collection.
func (p *Personal) One(ctx context.Context, id string) (*entity.PersonalData, error) {
	usr, err := p.DB.One(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "could not select one personal data")
	}
	return usr, nil
}

// All returns an array of personal information.
func (p *Personal) All(ctx context.Context) ([]*entity.PersonalData, error) {
	usrs, err := p.DB.All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "could not select all personal data")
	}
	return usrs, nil
}

// Update changes information in collection.
func (p *Personal) Update(ctx context.Context, document *entity.PersonalData) (int64, error) {
	return p.DB.Update(ctx, document)
}

// Remove deletes information from collection.
func (p *Personal) Remove(ctx context.Context, id string) (int64, error) {
	return p.DB.Remove(ctx, id)
}
