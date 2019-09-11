package mongo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/EreminDm/golang_basic_crud/entity"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRecive(t *testing.T) {
	oid := primitive.NewObjectID()

	tt := []struct {
		name      string
		enterT    entity.PersonalData
		expectedT personalData
		err       error
	}{
		{
			name: "Recive data from entity package type to mongo",
			enterT: entity.PersonalData{
				DocumentID:  oid.Hex(),
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
			expectedT: personalData{
				DocumentID:  &oid,
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
			// Using the variable on range scope `tc` in function literal (scopelint)
			actualT, err := receive(&tc.enterT)
			if tc.err != nil {
				assert.Equal(
					t,
					tc.err,
					err,
					fmt.Sprintf("errors not equal; want %v\n got: %v", tc.err, err),
				)
			}
			assert.NoError(
				t,
				err,
				fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err),
			)
			assert.Equal(
				t,
				tc.expectedT,
				actualT,
				fmt.Sprintf("expected type %v, actual %v", tc.expectedT, actualT),
			)
		})
	}
}

func TestTransmit(t *testing.T) {
	oid := primitive.NewObjectID()

	tt := []struct {
		name      string
		enterT    personalData
		expectedT entity.PersonalData
	}{
		{
			name: "transmit data from mongo package type to entity",
			enterT: personalData{
				DocumentID:  &oid,
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
			expectedT: entity.PersonalData{
				DocumentID:  oid.Hex(),
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
			actualT := tc.enterT.transmit()
			assert.Equal(
				t,
				tc.expectedT,
				actualT,
				fmt.Sprintf("expected type %v, actual %v", tc.expectedT, actualT),
			)
		})
	}
}

func TestInsert(t *testing.T) {
	oid := primitive.NewObjectID().Hex()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	m, err := Connect(ctx, "localhost:27017", collectionName)
	assert.NoError(t, err, "could not connect to db")
	tt := []struct {
		name       string
		collection *Mongodatabase
		enterT     entity.PersonalData
		ctx        context.Context
		err        error
	}{
		{
			name:       "",
			collection: m,
			enterT: entity.PersonalData{
				DocumentID:  oid,
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
			ctx: ctx,
			err: nil,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err = tc.collection.Insert(tc.ctx, &tc.enterT)
			if tc.err != nil {
				assert.Equal(
					t,
					tc.err,
					err,
					fmt.Sprintf("errors not equal; want %v\n got: %v", tc.err, err),
				)
				return
			}
			assert.NoError(t, err, "could not insert data to database")
			_, err = tc.collection.Remove(tc.ctx, tc.enterT.DocumentID)
			assert.NoError(t, err, "could not remove document from database")
		})
	}
}

func TestAll(t *testing.T) {
	oid := primitive.NewObjectID().Hex()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	m, err := Connect(ctx, "localhost:27017", collectionName)
	assert.NoError(t, err, "could not connect to db")
	tt := []struct {
		name       string
		collection *Mongodatabase
		expectedT  entity.PersonalData
		enterT     entity.PersonalData
		ctx        context.Context
		err        error
	}{
		{
			name:       "Select all without errors",
			collection: m,
			expectedT:  entity.PersonalData{},
			enterT: entity.PersonalData{
				DocumentID:  oid,
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
			ctx: ctx,
			err: nil,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err = tc.collection.Insert(tc.ctx, &tc.enterT)
			assert.NoError(t, err, "could not insert data to database")
			actualSlice, err := tc.collection.All(tc.ctx)
			if tc.err != nil {
				assert.Equal(
					t,
					tc.err,
					err,
					fmt.Sprintf("errors not equal; want %v\n got: %v", tc.err, err),
				)
			}
			assert.NoError(t, err, "could not select data from database")
			for _, aep := range actualSlice {
				assert.Equal(
					t,
					tc.err,
					err,
					fmt.Sprintf("actual data not equals; want %v\n got: %v", tc.expectedT, aep),
				)
			}
			_, err = tc.collection.Remove(tc.ctx, tc.enterT.DocumentID)
			assert.NoError(t, err, "could not remove document from database")
		})
	}
}

func TestOne(t *testing.T) {
	oid := primitive.NewObjectID().Hex()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	m, err := Connect(ctx, "localhost:27017", collectionName)
	assert.NoError(t, err, "could not connect to db")
	tt := []struct {
		name       string
		collection *Mongodatabase
		expectedT  entity.PersonalData
		enterT     entity.PersonalData
		ctx        context.Context
		err        error
	}{
		{
			name:       "Select one testing",
			collection: m,
			expectedT:  entity.PersonalData{},
			enterT: entity.PersonalData{
				DocumentID:  oid,
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
			ctx: ctx,
			err: nil,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err = tc.collection.Insert(tc.ctx, &tc.enterT)
			assert.NoError(t, err, "could not insert data to database")
			actualSlice, err := tc.collection.One(tc.ctx, tc.enterT.DocumentID)
			if tc.err != nil {
				assert.Equal(
					t,
					tc.err,
					err,
					fmt.Sprintf("errors not equal; want %v\n got: %v", tc.err, err),
				)
			}
			assert.NoError(t, err, "could not select data from database")

			assert.Equal(
				t,
				tc.err,
				err,
				fmt.Sprintf("actual data not equals; want %v\n got: %v", tc.expectedT, actualSlice),
			)

			_, err = tc.collection.Remove(tc.ctx, tc.enterT.DocumentID)
			assert.NoError(t, err, "could not remove document from database")
		})
	}
}

func TestRemove(t *testing.T) {
	oid := primitive.NewObjectID().Hex()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	m, err := Connect(ctx, "localhost:27017", collectionName)
	assert.NoError(t, err, "could not connect to db")
	tt := []struct {
		name             string
		collection       *Mongodatabase
		expectedResponce int
		enterT           entity.PersonalData
		ctx              context.Context
		err              error
	}{
		{
			name:             "Remove document from database",
			collection:       m,
			expectedResponce: 1,
			enterT: entity.PersonalData{
				DocumentID:  oid,
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
			ctx: ctx,
			err: nil,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err = tc.collection.Insert(tc.ctx, &tc.enterT)
			assert.NoError(t, err, "could not insert data to database")
			er, err := tc.collection.Remove(tc.ctx, tc.enterT.DocumentID)
			if tc.err != nil {
				assert.Equal(
					t,
					tc.err,
					err,
					fmt.Sprintf("errors not equal; want %v\n got: %v", tc.err, err),
				)
			}
			assert.NoError(t, err, "could not remove data from database")

			assert.Equal(
				t,
				tc.err,
				err,
				fmt.Sprintf("actual data not equals; want %v\n got: %v", tc.expectedResponce, er),
			)

			_, err = tc.collection.Remove(tc.ctx, tc.enterT.DocumentID)
			assert.NoError(t, err, "could not remove document from database")
		})
	}
}

func TestUpdate(t *testing.T) {
	oid := primitive.NewObjectID().Hex()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	m, err := Connect(ctx, "localhost:27017", collectionName)
	assert.NoError(t, err, "could not connect to db")
	tt := []struct {
		name             string
		collection       *Mongodatabase
		enterT           entity.PersonalData
		updateT          entity.PersonalData
		expectedResponce int
		ctx              context.Context
		err              error
	}{
		{
			name:       "Update Document in database",
			collection: m,
			enterT: entity.PersonalData{
				DocumentID:  oid,
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
			updateT: entity.PersonalData{
				DocumentID:  oid,
				Name:        "FirstName",
				LastName:    "LastName",
				Phone:       "999999999",
				Email:       "test@test.test",
				YearOfBirth: 1999,
			},
			expectedResponce: 1,
			ctx:              ctx,
			err:              nil,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err = tc.collection.Insert(tc.ctx, &tc.enterT)
			assert.NoError(t, err, "could not insert data to database")
			_, err := tc.collection.Update(tc.ctx, &tc.updateT)
			if tc.err != nil {
				assert.Equal(
					t,
					tc.err,
					err,
					fmt.Sprintf("errors not equal; want %v\n got: %v", tc.err, err),
				)
				return
			}
			assert.NoError(t, err, "could not insert data to database")

			_, err = tc.collection.Remove(tc.ctx, tc.enterT.DocumentID)
			assert.NoError(t, err, "could not remove document from database")
		})
	}
}
