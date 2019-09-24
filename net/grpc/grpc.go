package grpc

import (
	"context"
	n "net"

	"github.com/EreminDm/golang_basic_crud/controller"
	"github.com/EreminDm/golang_basic_crud/entity"
	"github.com/EreminDm/golang_basic_crud/net"
	grpcproto "github.com/EreminDm/golang_basic_crud/net/grpc/proto"
	"github.com/pkg/errors"
	grpc "google.golang.org/grpc"
)

// New returns implemented controller.
func New(c net.Provider) Controller {
	return Controller{
		CTR: c,
	}
}

// Controller describes controller implementation.
type Controller struct {
	CTR net.Provider
}

// ConnectServer runs on port 8888,
func ConnectServer(cp *controller.Personal) (n.Listener, *grpc.Server, error) {
	srv := grpc.NewServer()
	var pdServer = New(cp)
	grpcproto.RegisterPersonalDataServer(srv, pdServer)
	l, err := n.Listen("tcp", ":8888")
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not listen to :8888")
	}
	return l, srv, nil
}

// receive returns grpc package person object construction.
func receive(ep entity.PersonalData) grpcproto.Person {
	return grpcproto.Person{
		DocumentID:  ep.DocumentID,
		Name:        ep.Name,
		LastName:    ep.LastName,
		Phone:       ep.Phone,
		Email:       ep.Email,
		YearOfBirth: int32(ep.YearOfBirth),
	}
}

// transmit returns entity object construction.
func transmit(p *grpcproto.Person) entity.PersonalData {
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
func (c Controller) List(ctx context.Context, void *grpcproto.Void) (*grpcproto.PersonalDataList, error) {
	epList, err := c.CTR.All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "could not get personal information")
	}
	var pList grpcproto.PersonalDataList
	for _, ep := range epList {
		p := receive(ep)
		pList.Person = append(pList.Person, &p)

	}
	return &pList, nil
}

// One returns a person information by object id.
func (c Controller) One(ctx context.Context, documentID *grpcproto.ObjectID) (*grpcproto.Person, error) {
	ep, err := c.CTR.One(ctx, documentID.GetObjectID())
	if err != nil {
		return nil, errors.Wrap(err, "could not get personal information")
	}
	p := receive(ep)
	return &p, nil
}

// Update returns a count of updated documents.
func (c Controller) Update(ctx context.Context, personal *grpcproto.Person) (*grpcproto.Count, error) {
	ep := transmit(personal)
	count, err := c.CTR.Update(ctx, ep)
	if err != nil {
		return nil, errors.Wrap(err, "could not update document")
	}
	return &grpcproto.Count{
		Count: count,
	}, nil
}

// Insert returns added personal information objects.
func (c Controller) Insert(ctx context.Context, personal *grpcproto.Person) (*grpcproto.Person, error) {
	ep := transmit(personal)
	res, err := c.CTR.Insert(ctx, ep)
	if err != nil {
		return nil, errors.Wrap(err, "could not insert document")
	}
	p := receive(res)
	return &p, nil
}

// Remove returns a count of removed documents.
func (c Controller) Remove(ctx context.Context, documentID *grpcproto.ObjectID) (*grpcproto.Count, error) {
	oid := documentID.GetObjectID()
	count, err := c.CTR.Remove(ctx, oid)
	if err != nil {
		return nil, errors.Wrap(err, "could not remove document")
	}
	return &grpcproto.Count{
		Count: count,
	}, nil
}
