package grpc

import (
	context "context"

	"github.com/EreminDm/golang_basic_crud/entity"
	"github.com/EreminDm/golang_basic_crud/nets"
	"github.com/pkg/errors"
)

// New returns implemented controller.
func New(c nets.Provider) Controller {
	return Controller{
		CTR: c,
	}
}

// Controller describes controller implementation.
type Controller struct {
	CTR nets.Provider
}

// receive returns grpc package person object construction.
func receive(ep entity.PersonalData) Person {
	return Person{
		DocumentID:  ep.DocumentID,
		Name:        ep.Name,
		LastName:    ep.LastName,
		Phone:       ep.Phone,
		Email:       ep.Email,
		YearOfBirth: int32(ep.YearOfBirth),
	}
}

// transmit returns entity object construction.
func (p *Person) transmit() entity.PersonalData {
	return entity.PersonalData{
		DocumentID:  p.GetDocumentID(),
		Name:        p.GetName(),
		LastName:    p.GetLastName(),
		Phone:       p.GetPhone(),
		Email:       p.GetEmail(),
		YearOfBirth: int(p.GetYearOfBirth()),
	}
}

// List returns full personal data list.
func (c Controller) List(ctx context.Context, void *Void) (*PersonalDataList, error) {
	epList, err := c.CTR.All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "could not get personal information")
	}
	var pList PersonalDataList
	for _, ep := range epList {
		p := receive(ep)
		pList.Person = append(pList.Person, &p)

	}
	return &pList, nil
}

// One returns a person information by object id.
func (c Controller) One(ctx context.Context, documentID *ObjectID) (*Person, error) {
	ep, err := c.CTR.One(ctx, documentID.GetObjectID())
	if err != nil {
		return nil, errors.Wrap(err, "could not get personal information")
	}
	p := receive(ep)
	return &p, nil
}

// Update returns a count of updated documents.
func (c Controller) Update(ctx context.Context, personal *Person) (*Count, error) {
	ep := personal.transmit()
	count, err := c.CTR.Update(ctx, ep)
	if err != nil {
		return nil, errors.Wrap(err, "could not update document")
	}
	return &Count{
		Count: count,
	}, nil
}

// Insert returns added a personal information object.
func (c Controller) Insert(ctx context.Context, personal *Person) (*Person, error) {
	ep := personal.transmit()
	res, err := c.CTR.Insert(ctx, ep)
	if err != nil {
		return nil, errors.Wrap(err, "could not insert document")
	}
	p := receive(res)
	return &p, nil
}

// Remove returns a count of removed documents.
func (c Controller) Remove(ctx context.Context, documentID *ObjectID) (*Count, error) {
	oid := documentID.GetObjectID()
	count, err := c.CTR.Remove(ctx, oid)
	if err != nil {
		return nil, errors.Wrap(err, "could not remove document")
	}
	return &Count{
		Count: count,
	}, nil
}
