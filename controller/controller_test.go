package controller_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/EreminDm/golang_basic_crud/controller"
	"github.com/EreminDm/golang_basic_crud/entity"
	"github.com/EreminDm/golang_basic_crud/mongo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNew(t *testing.T) {
	var expected = &controller.Personal{}
	var cp controller.DBProvider

	tt := []struct {
		name     string
		provider controller.DBProvider
		equal    bool
	}{
		{name: "New controller", provider: cp, equal: true},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			actual := controller.New(tc.provider)
			if tc.equal != assert.IsType(t, expected, actual) {
				t.Fatalf("not equals interfaces, expected: %v, actual: %v", expected, actual)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	m, err := mongo.Connect(ctx, "localhost:27017", "test")
	if err != nil {
		log.Fatalf(`couldn't connect to database: %v`, err)
	}
	// returns controller provider.
	c := controller.New(m)

	defer cancel()
	tt := []struct {
		name     string
		provider controller.DBProvider
		context  context.Context
		document entity.PersonalData
	}{
		{
			name:     "Insert controller",
			provider: c,
			context:  ctx,
			document: entity.PersonalData{
				DocumentID:  primitive.NewObjectID().Hex(),
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.provider.Insert(tc.context, &tc.document)
			assert.NoError(t, err, "could not insert data")
			_, err = tc.provider.Remove(tc.context, actual.DocumentID)
			assert.NoError(t, err, "could not remove data")
		})
	}
}
