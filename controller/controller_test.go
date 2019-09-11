package controller_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/EreminDm/golang_basic_crud/controller"
	"github.com/EreminDm/golang_basic_crud/entity"
	"github.com/EreminDm/golang_basic_crud/mongo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type controllerMockedObject struct {
	mock.Mock
}

func (m *controllerMockedObject) Insert(ctx context.Context, document *entity.PersonalData) (entity.PersonalData, error) {
	fmt.Println("Mocked insert function")
	fmt.Printf("Document passed in: %v\n", document)
	ep, err := m.Insert(ctx, document)

	return ep, errors.Wrap(err, "could not insert personal data")
}

func (m *controllerMockedObject) DummyFunc() {
	fmt.Println("Dummy")
}

func (m *controllerMockedObject) One(ctx context.Context, id string) (entity.PersonalData, error) {
	fmt.Println("Mocked one function")
	fmt.Printf("ID passed in: %s\n", id)
	usr, err := m.One(ctx, id)
	return usr, errors.Wrap(err, "could not select one personal data")

}

// All returns an array of personal information.
func (m *controllerMockedObject) All(ctx context.Context) ([]entity.PersonalData, error) {
	fmt.Println("Mocked all function")
	usrs, err := m.All(ctx)
	return usrs, errors.Wrap(err, "could not select all personal data")
}

// Update changes information in collection.
func (m *controllerMockedObject) Update(ctx context.Context, document *entity.PersonalData) (int64, error) {
	fmt.Println("Mocked update function")
	fmt.Printf("Document passed in: %v\n", document)
	count, err := m.Update(ctx, document)
	return count, errors.Wrap(err, "could not update personal data")
}

// Remove deletes information from collection.
func (m *controllerMockedObject) Remove(ctx context.Context, id string) (int64, error) {
	fmt.Println("Mocked remove function")
	fmt.Printf("ID passed in: %s\n", id)
	count, err := m.Remove(ctx, id)
	return count, errors.Wrap(err, "could not delete personal data")
}

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
	defer cancel()
	c := new(controllerMockedObject)

	tt := []struct {
		name     string
		context  context.Context
		document entity.PersonalData
	}{
		{
			name:    "Insert controller",
			context: ctx,
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
			c.On("Insert", tc.context, &tc.document).Return(true)
			_, err := c.Insert(tc.context, &tc.document)
			assert.NoError(t, err, "could not insert data")
			_, err = c.Remove(tc.context, tc.document.DocumentID)
			assert.NoError(t, err, "could not remove data")

			c.AssertExpectations(t)
		})
	}
}

func TestOne(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	var oid = primitive.NewObjectID().Hex()
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
		expected entity.PersonalData
	}{
		{
			name:     "One controller",
			provider: c,
			context:  ctx,
			document: entity.PersonalData{
				DocumentID:  oid,
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
			expected: entity.PersonalData{
				DocumentID:  oid,
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
			_, err := tc.provider.Insert(tc.context, &tc.document)
			assert.NoError(t, err, "could not insert data")
			actual, err := tc.provider.One(tc.context, oid)
			assert.NoError(t, err, "could not select data")
			assert.Equal(t,
				tc.expected,
				actual,
				fmt.Sprintf("expected object %v is not equals %v", tc.expected, actual),
			)
			_, err = tc.provider.Remove(tc.context, tc.document.DocumentID)
			assert.NoError(t, err, "could not remove data")
		})
	}
}

func TestAll(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	var oid = primitive.NewObjectID().Hex()
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
		expected []entity.PersonalData
	}{
		{
			name:     "All controller",
			provider: c,
			context:  ctx,
			document: entity.PersonalData{
				DocumentID:  oid,
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
			expected: []entity.PersonalData{
				{
					DocumentID:  oid,
					Name:        "Name",
					LastName:    "LName",
					Phone:       "1235486",
					Email:       "test@test.test",
					YearOfBirth: 1234,
				},
			},
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.provider.Insert(tc.context, &tc.document)
			assert.NoError(t, err, "could not insert data")
			actual, err := tc.provider.All(tc.context)
			assert.NoError(t, err, "could not select data")
			assert.IsType(t,
				tc.expected,
				actual,
				fmt.Sprintf("expected object %v is not equals %v", tc.expected, actual),
			)
			_, err = tc.provider.Remove(tc.context, tc.document.DocumentID)
			assert.NoError(t, err, "could not remove data")
		})
	}
}

func TestUpdate(t *testing.T) {
	oid := primitive.NewObjectID().Hex()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	m, err := mongo.Connect(ctx, "localhost:27017", "test")
	if err != nil {
		log.Fatalf(`couldn't connect to database: %v`, err)
	}
	// returns controller provider.
	c := controller.New(m)

	defer cancel()
	tt := []struct {
		name         string
		provider     controller.DBProvider
		context      context.Context
		document     entity.PersonalData
		newDocument  entity.PersonalData
		updatesCount int64
	}{
		{
			name:     "Update controller",
			provider: c,
			context:  ctx,
			document: entity.PersonalData{
				DocumentID:  oid,
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
			newDocument: entity.PersonalData{
				DocumentID:  oid,
				Name:        "FirstName",
				LastName:    "LastName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
			updatesCount: 1,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.provider.Insert(tc.context, &tc.document)
			assert.NoError(t, err, "could not insert data")
			k, err := tc.provider.Update(tc.context, &tc.newDocument)
			assert.NoError(t, err, "could not update data")
			assert.Equal(t, tc.updatesCount, k, "not equals")
			_, err = tc.provider.Remove(tc.context, tc.document.DocumentID)
			assert.NoError(t, err, "could not remove data")
		})
	}
}
