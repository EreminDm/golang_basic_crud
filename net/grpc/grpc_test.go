package grpc

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/EreminDm/golang_basic_crud/controller"
	"github.com/EreminDm/golang_basic_crud/entity"
	"github.com/EreminDm/golang_basic_crud/net"
	grpcproto "github.com/EreminDm/golang_basic_crud/net/grpc/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNew(t *testing.T) {
	var p net.Provider
	var expected Controller

	tt := []struct {
		name     string
		provider net.Provider
		equal    bool
	}{
		{name: "Not nil interface", provider: p, equal: true},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			actual := New(p)
			if tc.equal != assert.IsType(t, expected, actual) {
				t.Fatalf("not equals interfaces, expected: %v, actual: %v", expected, actual)
			}
		})
	}
}

func TestConnectServer(t *testing.T) {
	c := &controller.Personal{}
	_, _, err := ConnectServer(c)
	assert.NoError(t, err, "could not make correct connection")
	//assert.Equal(t, n.Listener, l, "type is not equals want %v, got %v", n.Listener, l)
}

func TestReceive(t *testing.T) {
	tt := []struct {
		name      string
		enterT    entity.PersonalData
		expectedT grpcproto.Person
		err       error
	}{
		{
			name: "Recive data from entity package type to mongo",
			enterT: entity.PersonalData{
				DocumentID:  "ObjectID",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
			expectedT: grpcproto.Person{
				DocumentID:  "ObjectID",
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
			actualT := receive(tc.enterT)
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
	tt := []struct {
		name      string
		enterT    grpcproto.Person
		expectedT entity.PersonalData
		err       error
	}{
		{
			name: "Recive data from entity package type to mongo",
			enterT: grpcproto.Person{
				DocumentID:  "ObjectID",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1234,
			},
			expectedT: entity.PersonalData{
				DocumentID:  "ObjectID",
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
			actualT := transmit(&tc.enterT)
			assert.Equal(
				t,
				tc.expectedT,
				actualT,
				fmt.Sprintf("expected type %v, actual %v", tc.expectedT, actualT),
			)
		})
	}
}

type controllerMockedObject struct {
	mock.Mock
}

// Insert adds a personal information.
func (m *controllerMockedObject) Insert(ctx context.Context, document entity.PersonalData) (entity.PersonalData, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(entity.PersonalData), args.Error(1)
}

// One returns personal information.
func (m *controllerMockedObject) One(ctx context.Context, id string) (entity.PersonalData, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.PersonalData), args.Error(1)
}

// All returns an array of personal information.
func (m *controllerMockedObject) All(ctx context.Context) ([]entity.PersonalData, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.PersonalData), args.Error(1)
}

// Update changes information in collection.
func (m *controllerMockedObject) Update(ctx context.Context, document entity.PersonalData) (int64, error) {
	args := m.Called(ctx, document)
	return int64(args.Int(0)), args.Error(1)
}

// Remove deletes information from collection.
func (m *controllerMockedObject) Remove(ctx context.Context, id string) (int64, error) {
	args := m.Called(ctx, id)
	return int64(1), args.Error(1)
}

func TestInsert(t *testing.T) {
	ctr := new(controllerMockedObject)
	c := New(ctr)
	tt := []struct {
		name           string
		object         grpcproto.Person
		expectedObject entity.PersonalData
		expectedError  error
		err            string
	}{
		{
			name: "Success request",
			object: grpcproto.Person{
				DocumentID:  "",
				Name:        "firstName",
				LastName:    "secondName",
				Phone:       "",
				Email:       "",
				YearOfBirth: 1980,
			},
			expectedError: nil,
		}, {
			name:          "Wrong request",
			object:        grpcproto.Person{},
			expectedError: errors.New("error"),
			err:           "could not insert document: error",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctr.On("Insert", mock.Anything, transmit(&tc.object)).Return(tc.expectedObject, tc.expectedError).Once()
			p, err := c.Insert(context.Background(), &tc.object)
			if tc.err != "" {
				assert.Equal(t,
					tc.err,
					err.Error(),
					"not equals",
				)
				return
			}
			assert.NoError(t, err, fmt.Sprintf("could not insert data: %v", err))
			assert.IsType(t,
				&tc.object,
				p,
				fmt.Sprintf("expected value %v; got %v", tc.object, p),
			)
		})
	}
}

func TestOne(t *testing.T) {
	ctr := new(controllerMockedObject)
	c := New(ctr)
	oid := primitive.NewObjectID().Hex()

	tt := []struct {
		name           string
		object         grpcproto.Person
		objID          *grpcproto.ObjectID
		expectedObject entity.PersonalData
		expectedError  error
		err            string
	}{
		{
			name: "Success request",
			objID: &grpcproto.ObjectID{
				ObjectID: oid,
			},
			expectedError: nil,
			err:           "",
		}, {
			name:          "Wrong request",
			objID:         &grpcproto.ObjectID{},
			expectedError: errors.New("error"),
			err:           "could not get personal information: error",
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctr.On("One", mock.Anything, tc.objID.GetObjectID()).Return(tc.expectedObject, tc.expectedError)
			p, err := c.One(context.Background(), tc.objID)
			if tc.err != "" {
				assert.Equal(t,
					tc.err,
					err.Error(),
					"not equals",
				)
				return
			}
			assert.NoError(t, err, fmt.Sprintf("could not return data: %v", err))
			assert.IsType(t,
				&tc.object,
				p,
				fmt.Sprintf("expected value %v; got %v", tc.object, p),
			)
		})
	}
}

func TestList(t *testing.T) {
	ctr := new(controllerMockedObject)
	c := New(ctr)

	tt := []struct {
		name           string
		object         grpcproto.PersonalDataList
		expectedObject []entity.PersonalData
		expectedError  error
		err            string
	}{
		{
			name:          "Success request",
			expectedError: nil,
			err:           "",
		},
		// {
		// 	name:          "Wrong request",
		// 	expectedError: fmt.Errorf("error"),
		// 	err:           "could not get personal information: error",
		// },
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctr.On("All", mock.Anything).Return(tc.expectedObject, tc.expectedError)
			p, err := c.List(context.Background(), &grpcproto.Void{})
			if tc.err != "" {
				assert.Equal(t,
					tc.err,
					err,
					"not equals",
				)
				return
			}
			assert.NoError(t, err, fmt.Sprintf("could not return data: %v", err))
			assert.IsType(t,
				&tc.object,
				p,
				fmt.Sprintf("expected value %v; got %v", tc.object, p),
			)
		})
	}
}

func TestUpdate(t *testing.T) {
	ctr := new(controllerMockedObject)
	c := New(ctr)
	tt := []struct {
		name          string
		object        grpcproto.Person
		expectedCount grpcproto.Count
		expectedError error
		err           string
	}{
		{
			name: "Success request",
			object: grpcproto.Person{
				DocumentID:  "",
				Name:        "firstName",
				LastName:    "secondName",
				Phone:       "",
				Email:       "",
				YearOfBirth: 1980,
			},
			expectedCount: grpcproto.Count{
				Count: 1,
			},
			expectedError: nil,
		},
		{
			name:   "Wrong request",
			object: grpcproto.Person{},
			expectedCount: grpcproto.Count{
				Count: 0,
			},
			expectedError: errors.New("error"),
			err:           "could not update document: error",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctr.On("Update", mock.Anything, transmit(&tc.object)).Return(int(tc.expectedCount.GetCount()), tc.expectedError).Once()
			p, err := c.Update(context.Background(), &tc.object)
			if tc.err != "" {
				assert.Equal(t,
					tc.err,
					err.Error(),
					"not equals",
				)
				return
			}
			assert.NoError(t, err, fmt.Sprintf("could not insert data: %v", err))
			assert.Equal(t,
				&tc.expectedCount,
				p,
				fmt.Sprintf("expected value %v; got %v", tc.object, p),
			)
		})
	}
}

func TestRemove(t *testing.T) {
	ctr := new(controllerMockedObject)
	c := New(ctr)
	oid := primitive.NewObjectID().Hex()

	tt := []struct {
		name          string
		object        grpcproto.Person
		objID         *grpcproto.ObjectID
		expectedCount grpcproto.Count
		expectedError error
		err           string
	}{
		{
			name: "Success request",
			objID: &grpcproto.ObjectID{
				ObjectID: oid,
			},
			expectedCount: grpcproto.Count{
				Count: 1,
			},
			expectedError: nil,
			err:           "",
		}, {
			name:  "Wrong request",
			objID: &grpcproto.ObjectID{},
			expectedCount: grpcproto.Count{
				Count: 1,
			},
			expectedError: errors.New("error"),
			err:           "could not remove document: error",
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctr.On("Remove", mock.Anything, tc.objID.GetObjectID()).Return(int(tc.expectedCount.GetCount()), tc.expectedError)
			p, err := c.Remove(context.Background(), tc.objID)
			if tc.err != "" {
				assert.Equal(t,
					tc.err,
					err.Error(),
					"not equals",
				)
				return
			}
			assert.NoError(t, err, fmt.Sprintf("could not return data: %v", err))
			assert.Equal(t,
				&tc.expectedCount,
				p,
				fmt.Sprintf("expected value %v; got %v", tc.object, p),
			)
		})
	}
}
