package controller_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/EreminDm/golang_basic_crud/controller"
	"github.com/EreminDm/golang_basic_crud/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DBMockedObject struct {
	mock.Mock
}

func (m *DBMockedObject) Insert(ctx context.Context, document entity.PersonalData) (entity.PersonalData, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(entity.PersonalData), args.Error(1)
}

func (m *DBMockedObject) One(ctx context.Context, id string) (entity.PersonalData, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.PersonalData), args.Error(1)
}

// All returns an array of personal information.
func (m *DBMockedObject) All(ctx context.Context) ([]entity.PersonalData, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.PersonalData), args.Error(1)
}

// Update changes information in collection.
func (m *DBMockedObject) Update(ctx context.Context, document entity.PersonalData) (int64, error) {
	args := m.Called(ctx, document)
	return int64(args.Int(0)), args.Error(1)
}

// Remove deletes information from collection.
func (m *DBMockedObject) Remove(ctx context.Context, id string) (int64, error) {
	args := m.Called(ctx, id)
	return int64(1), args.Error(1)
}

func TestNew(t *testing.T) {
	var expected = &controller.Personal{}
	var cp controller.DBProvider

	tt := []struct {
		name     string
		provider controller.DBProvider
		equal    bool
	}{
		{name: "New midleware controller", provider: cp, equal: true},
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
	db := new(DBMockedObject)
	c := controller.New(db)
	tt := []struct {
		name     string
		context  context.Context
		document entity.PersonalData
		err      error
	}{
		{
			name:    "Controller -> Insert",
			context: ctx,
			document: entity.PersonalData{
				DocumentID:  primitive.NewObjectID().Hex(),
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
			err: nil,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			db.On("Insert", tc.context, tc.document).Return(tc.document, tc.err).Once()
			_, err := c.Insert(tc.context, tc.document)
			assert.NoError(t, err, "could not insert data")
			db.AssertExpectations(t)
		})
	}
}

func TestOne(t *testing.T) {
	var oid = primitive.NewObjectID().Hex()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db := new(DBMockedObject)
	c := controller.New(db)
	tt := []struct {
		name     string
		context  context.Context
		expected entity.PersonalData
		err      error
	}{
		{
			name:    "Controller -> One",
			context: ctx,
			expected: entity.PersonalData{
				DocumentID:  oid,
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
			err: nil,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			db.On("One", tc.context, oid).Return(tc.expected, tc.err).Once()
			actual, err := c.One(tc.context, oid)
			assert.NoError(t, err, "could not select data")
			assert.Equal(t,
				tc.expected,
				actual,
				fmt.Sprintf("expected object %v is not equals %v", tc.expected, actual),
			)
			db.AssertExpectations(t)
		})
	}
}

func TestAll(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var oid = primitive.NewObjectID().Hex()
	db := new(DBMockedObject)
	c := controller.New(db)

	tt := []struct {
		name     string
		context  context.Context
		expected []entity.PersonalData
		err      error
	}{
		{
			name:    "Controller -> All",
			context: ctx,
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
			err: nil,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			db.On("All", tc.context).Return(tc.expected, tc.err).Once()
			actual, err := c.All(tc.context)
			assert.NoError(t, err, "could not select data")
			assert.IsType(t,
				tc.expected,
				actual,
				fmt.Sprintf("expected object %v is not equals %v", tc.expected, actual),
			)
			db.AssertExpectations(t)
		})
	}
}

func TestUpdate(t *testing.T) {
	oid := primitive.NewObjectID().Hex()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db := new(DBMockedObject)
	c := controller.New(db)

	tt := []struct {
		name         string
		context      context.Context
		newDocument  entity.PersonalData
		updatesCount int64
		err          error
	}{
		{
			name:    "Controller -> Update",
			context: ctx,
			newDocument: entity.PersonalData{
				DocumentID:  oid,
				Name:        "FirstName",
				LastName:    "LastName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
			updatesCount: 1,
			err:          nil,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			db.On("Update", tc.context, tc.newDocument).Return(1, tc.err).Once()
			k, err := c.Update(tc.context, tc.newDocument)
			assert.NoError(t, err, "could not update data")
			assert.Equal(t, tc.updatesCount, k, "not equals")
			db.AssertExpectations(t)
		})
	}
}

func TestRemove(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db := new(DBMockedObject)
	c := controller.New(db)

	tt := []struct {
		name        string
		context     context.Context
		oid         string
		removeCount int64
		err         error
	}{
		{
			name:        "Controller -> Remove",
			context:     ctx,
			oid:         primitive.NewObjectID().Hex(),
			removeCount: 1,
			err:         nil,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			db.On("Remove", tc.context, tc.oid).Return(1, tc.err).Once()
			k, err := c.Remove(tc.context, tc.oid)
			assert.NoError(t, err, "could not remove data")
			assert.Equal(t, tc.removeCount, k, "not equals")
			db.AssertExpectations(t)
		})
	}
}
