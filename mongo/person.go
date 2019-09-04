package mongo

import (
	"context"

	"github.com/EreminDm/golang_basic_crud/entity"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// personalData description.
type personalData struct {
	DocumentID  *primitive.ObjectID `bson:"_id"`
	Name        string              `bson:"name"`
	LastName    string              `bson:"lastName"`
	Phone       string              `bson:"phone,omitempty"`
	Email       string              `bson:"email,omitempty"`
	YearOfBirth int                 `bson:"yaerOfBirth,omitempty"`
}

// receive returns mongo personal data construction.
func receive(ep *entity.PersonalData) (*personalData, error) {

	// returns ObjectID type from string
	oid, err := primitive.ObjectIDFromHex(ep.DocumentID)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert DocumentID type string to type ObjectID")
	}
	return &personalData{
		DocumentID:  &oid,
		Name:        ep.Name,
		LastName:    ep.LastName,
		Phone:       ep.Phone,
		Email:       ep.Email,
		YearOfBirth: ep.YearOfBirth,
	}, nil
}

// transmit returns entity data construction.
func (p *personalData) transmit() *entity.PersonalData {
	return &entity.PersonalData{
		DocumentID:  p.DocumentID.Hex(),
		Name:        p.Name,
		LastName:    p.LastName,
		Phone:       p.Phone,
		Email:       p.Email,
		YearOfBirth: p.YearOfBirth,
	}
}

// One returns personal data for a given id,
// id params to make filtration.
func (m *Mongodatabase) One(ctx context.Context, id string) (*entity.PersonalData, error) {
	var p personalData
	val, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't decode object id from hex err")
	}
	filter := bson.D{{"_id", val}}
	err = m.Person.FindOne(ctx, filter).Decode(&p)
	if err != nil {
		return nil, errors.Wrap(err, "could not find collection ")
	}
	return p.transmit(), nil
}

// Insert is a function which adding data to database.
func (m *Mongodatabase) Insert(ctx context.Context, document *entity.PersonalData) (interface{}, error) {
	p, err := receive(document)
	if err != nil {
		return nil, errors.Wrap(err, "could not receive data")
	}
	result, err := m.Person.InsertOne(ctx, p)
	if err != nil {
		return nil, errors.Wrap(err, "could not add document(s) in database")
	}
	return result.InsertedID, nil
}

// All selects all documents from database.
func (m *Mongodatabase) All(ctx context.Context) ([]*entity.PersonalData, error) {
	// no filter by default.
	// Searches documents in colletion.
	cursor, err := m.Person.Find(ctx, nil, options.Find())
	if err != nil {
		return nil, errors.Wrap(err, "could not find document in database")
	}
	defer cursor.Close(ctx)
	// Decode documents from colletion.
	var pArr []personalData
	err = cursor.All(ctx, &pArr)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode document to struct")
	}
	if err = cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "curser error")
	}
	// converting structs
	var epArr []*entity.PersonalData
	for _, p := range pArr {
		epArr = append(epArr, p.transmit())
	}
	return epArr, nil
}

// Remove deletes document from Mongo.
func (m *Mongodatabase) Remove(ctx context.Context, id string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, errors.Wrap(err, "couldn't decode object id from hex err")
	}
	filter := bson.D{{"_id", objectID}}
	delResult, err := m.Person.DeleteOne(ctx, filter)
	if err != nil {
		return 0, errors.Wrap(err, "could not remove document")
	}

	return delResult.DeletedCount, nil
}

// Update rewrites information in db by user id filtration.
func (m *Mongodatabase) Update(ctx context.Context, ep *entity.PersonalData) (int64, error) {
	p, err := receive(ep)
	if err != nil {
		return 0, errors.Wrap(err, "couldnt receive struct")
	}
	filter := bson.D{{"_id", p.DocumentID}}
	update := bson.D{{
		"$in", bson.D{{"name", p.Name}, {"lastName", p.LastName}, {"phone", p.Phone}, {"email", p.Email}, {"yaerOfBirth", p.YearOfBirth}},
	}}
	updateResult, err := m.Person.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, errors.Wrap(err, "could not update object")
	}
	return updateResult.ModifiedCount, nil
}
