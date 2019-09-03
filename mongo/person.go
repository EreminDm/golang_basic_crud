package mongo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PersonalData description.
type PersonalData struct {
	DocumentID  string `bson:"_id"` // as *primitive.ObjectID
	Name        string `bson:"name"`
	LastName    string `bson:"lastName"`
	Phone       string `bson:"phone,omitempty"`
	Email       string `bson:"email,omitempty"`
	YearOfBirth int    `bson:"yaerOfBirth,omitempty"`
}

// One returns personal data for a given id,
// id params to make filtration.
func (m *Mongodatabase) One(ctx context.Context, id string) (result *PersonalData, err error) {
	val, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't decode object id from hex err")
	}
	filter := bson.D{{"_id", val}}
	err = m.Person.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, errors.Wrap(err, "could not find collection ")
	}
	return result, nil
}

// Insert is a function which adding data to database.
func (m *Mongodatabase) Insert(ctx context.Context, document *PersonalData) (interface{}, error) {
	result, err := m.Person.InsertOne(ctx, document)
	if err != nil {
		return nil, errors.Wrap(err, "could not add document(s) in database")
	}
	return result.InsertedID, nil
}

// All selects all documents from database.
func (m *Mongodatabase) All(ctx context.Context) (results *[]PersonalData, err error) {
	// no filter by default.
	// Searches documents in colletion.
	cursor, err := m.Person.Find(ctx, nil, options.Find())
	if err != nil {
		return nil, errors.Wrap(err, "could not find document in database")
	}
	defer cursor.Close(ctx)
	// Decode documents from colletion.
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode document to struct")
	}
	if err = cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "curser error")
	}
	return results, nil
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
func (m *Mongodatabase) Update(ctx context.Context, p *PersonalData) (int64, error) {
	docID, err := primitive.ObjectIDFromHex(p.DocumentID)
	if err != nil {
		return 0, errors.Wrap(err, "could not make object id from incoming hex")
	}
	filter := bson.D{{"_id", docID}}
	update := bson.D{{
		"$in", bson.D{{"name", p.Name}, {"lastName", p.LastName}, {"phone", p.Phone}, {"email", p.Email}, {"yaerOfBirth", p.YearOfBirth}},
	}}
	updateResult, err := m.Person.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, errors.Wrap(err, "could not update object")
	}
	return updateResult.ModifiedCount, nil
}
