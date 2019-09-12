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
	One(ctx context.Context, value string) (entity.PersonalData, error)
	All(ctx context.Context) ([]entity.PersonalData, error)
	Remove(ctx context.Context, id string) (int64, error)
	Update(ctx context.Context, p entity.PersonalData) (int64, error)
	Insert(ctx context.Context, document entity.PersonalData) (entity.PersonalData, error)
}

// New returns new Personal provider.
func New(db DBProvider) *Personal {
	return &Personal{
		DB: db,
	}
}

// Insert adds data to collection.
func (p *Personal) Insert(ctx context.Context, document entity.PersonalData) (entity.PersonalData, error) {
	ep, err := p.DB.Insert(ctx, document)
	return ep, errors.Wrap(err, "could not insert personal data")
}

// One returns personal data from collection.
func (p *Personal) One(ctx context.Context, id string) (entity.PersonalData, error) {
	usr, err := p.DB.One(ctx, id)
	return usr, errors.Wrap(err, "could not select one personal data")

}

// All returns an array of personal information.
func (p *Personal) All(ctx context.Context) ([]entity.PersonalData, error) {
	usrs, err := p.DB.All(ctx)
	return usrs, errors.Wrap(err, "could not select all personal data")
}

// Update changes information in collection.
func (p *Personal) Update(ctx context.Context, document entity.PersonalData) (int64, error) {
	count, err := p.DB.Update(ctx, document)
	return count, errors.Wrap(err, "could not update personal data")
}

// Remove deletes information from collection.
func (p *Personal) Remove(ctx context.Context, id string) (int64, error) {
	count, err := p.DB.Remove(ctx, id)
	return count, errors.Wrap(err, "could not delete personal data")
}
