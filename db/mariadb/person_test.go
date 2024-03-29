package mariadb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/EreminDm/golang_basic_crud/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecive(t *testing.T) {
	tt := []struct {
		name      string
		enterT    entity.PersonalData
		expectedT personalData
	}{
		{
			name: "MariaDB -> Recive data from entity package type to mongo",
			enterT: entity.PersonalData{
				DocumentID:  "KGPFOFOGB",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1999,
			},
			expectedT: personalData{
				ID:          "KGPFOFOGB",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1999,
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Using the variable on range scope `tc` in function literal (scopelint)
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
		enterT    personalData
		expectedT entity.PersonalData
	}{
		{
			name: "MariaDB -> Transmit data from mongo package type to entity",
			enterT: personalData{
				ID:          "KGPFOFOGB",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1999,
			},
			expectedT: entity.PersonalData{
				DocumentID:  "KGPFOFOGB",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1999,
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	wrongCTX, wCanel := context.WithCancel(context.Background())
	wCanel()
	m, err := Connect(ctx, conURI, dbName, 1, 1)
	require.NoError(t, err, "could not connect to db")
	tt := []struct {
		name       string
		collection *MariaDB
		enterT     entity.PersonalData
		ctx        context.Context
		err        string
	}{
		{
			name:       "MariaDB -> Success Insert",
			collection: m,
			enterT: entity.PersonalData{
				DocumentID:  "KGPFOFOGB",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1999,
			},
			ctx: ctx,
			err: "",
		},
		{
			name:       "MariaDB -> Wrong Insert canceled context",
			collection: m,
			enterT: entity.PersonalData{
				DocumentID:  "KGPFOFOGB",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1999,
			},
			ctx: wrongCTX,
			err: "could not exec query statement: context canceled",
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err = tc.collection.Insert(tc.ctx, tc.enterT)
			if tc.err != "" {
				assert.Equal(
					t,
					tc.err,
					err.Error(),
					fmt.Sprintf("errors not equal; want %v\n got: %v", tc.err, err.Error()),
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
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	wrongCTX, cancelCTX := context.WithCancel(context.Background())
	cancelCTX()

	m, err := Connect(ctx, conURI, dbName, 1, 1)
	require.NoError(t, err, "could not connect to db")
	tt := []struct {
		name       string
		collection *MariaDB
		expectedT  entity.PersonalData
		enterT     entity.PersonalData
		ctx        context.Context
		err        string
	}{
		{
			name:       "MariaDB -> Select all without errors",
			collection: m,
			expectedT:  entity.PersonalData{},
			enterT: entity.PersonalData{
				DocumentID:  "KGPFOFOGB",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1999,
			},
			ctx: ctx,
			err: "",
		},
		{
			name:       "MariaDB -> Select all with wrong context",
			collection: m,
			expectedT:  entity.PersonalData{},
			enterT: entity.PersonalData{
				DocumentID:  "KGPFOFOGB",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1999,
			},
			ctx: wrongCTX,
			err: "could not make query: context canceled",
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err = tc.collection.Insert(ctx, tc.enterT)
			assert.NoError(t, err, "could not insert data to database")
			actualSlice, err := tc.collection.All(tc.ctx)
			if tc.err != "" {
				assert.Equal(
					t,
					tc.err,
					err.Error(),
					fmt.Sprintf("errors not equal; want %v\n got: %v", tc.err, err.Error()),
				)
				_, err = tc.collection.Remove(ctx, tc.enterT.DocumentID)
				assert.NoError(t, err, "could not remove document from database")
				return
			}
			assert.NoError(t, err, "could not select data from database")
			for _, aep := range actualSlice {
				assert.IsType(
					t,
					tc.expectedT,
					aep,
					fmt.Sprintf("actual data not equals; want %v\n got: %v", tc.expectedT, aep),
				)
			}
			_, err = tc.collection.Remove(tc.ctx, tc.enterT.DocumentID)
			assert.NoError(t, err, "could not remove document from database")
		})
	}
}

func TestOne(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	m, err := Connect(ctx, conURI, dbName, 1, 1)
	require.NoError(t, err, "could not connect to db")
	tt := []struct {
		name       string
		collection *MariaDB
		expectedT  entity.PersonalData
		enterT     entity.PersonalData
		oid        string
		ctx        context.Context
		err        string
	}{
		{
			name:       "MariaDB -> Select one testing",
			collection: m,
			expectedT:  entity.PersonalData{},
			enterT: entity.PersonalData{
				DocumentID:  "KGPFOFOGB",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1999,
			},
			oid: "KGPFOFOGB",
			ctx: ctx,
			err: "",
		},
		{
			name:       "MariaDB -> Wrong Select one testing",
			collection: m,
			expectedT:  entity.PersonalData{},
			enterT: entity.PersonalData{
				DocumentID:  "KGPFOFOGB",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1999,
			},
			oid: "abc",
			ctx: ctx,
			err: "could not scan row: sql: no rows in result set",
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err = tc.collection.Insert(ctx, tc.enterT)
			assert.NoError(t, err, "could not insert data to database")
			actualSlice, err := tc.collection.One(tc.ctx, tc.oid)
			if tc.err != "" {
				assert.Equal(
					t,
					tc.err,
					err.Error(),
					fmt.Sprintf("errors not equal; want %v\n got: %v", tc.err, err.Error()),
				)
				_, err = tc.collection.Remove(ctx, tc.enterT.DocumentID)
				assert.NoError(t, err, "could not remove document from database")
				return
			}
			assert.NoError(t, err, "could not select data from database")
			assert.IsType(
				t,
				tc.expectedT,
				actualSlice,
				fmt.Sprintf("actual data not equals; want %v\n got: %v", tc.expectedT, actualSlice),
			)
			_, err = tc.collection.Remove(ctx, tc.enterT.DocumentID)
			assert.NoError(t, err, "could not remove document from database")
		})
	}
}

func TestRemove(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	wrongctx, wrcancel := context.WithTimeout(context.Background(), 5*time.Second)
	wrcancel()
	m, err := Connect(ctx, conURI, dbName, 1, 1)
	require.NoError(t, err, "could not connect to db")
	tt := []struct {
		name             string
		collection       *MariaDB
		expectedResponce int64
		enterT           entity.PersonalData
		removingOID      string
		ctx              context.Context
		err              string
	}{
		{
			name:             "MariaDB -> Remove document from database",
			collection:       m,
			expectedResponce: 1,
			enterT: entity.PersonalData{
				DocumentID:  "KGPFOFOGB",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1999,
			},
			removingOID: "KGPFOFOGB",
			ctx:         ctx,
			err:         "",
		},
		{
			name:             "MariaDB -> Fail remove document from database",
			collection:       m,
			expectedResponce: 0,
			enterT: entity.PersonalData{
				DocumentID:  "KGPFOFOGB",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1999,
			},
			removingOID: "IJOGSDIOJGDSIOD",
			ctx:         ctx,
			err:         "",
		},
		{
			name:             "MariaDB -> wrong context -> failed remove document from database",
			collection:       m,
			expectedResponce: 0,
			enterT: entity.PersonalData{
				DocumentID:  "KGPFOFOGB",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1999,
			},
			removingOID: "KOIFIJFJIFSA",
			ctx:         wrongctx,
			err:         "could not remove data: context canceled",
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err = tc.collection.Insert(ctx, tc.enterT)
			assert.NoError(t, err, "could not insert data to database")
			er, err := tc.collection.Remove(tc.ctx, tc.removingOID)
			if tc.expectedResponce != er {
				assert.Equal(
					t,
					tc.expectedResponce,
					er,
					fmt.Sprintf("expected response not equal; want %v\n got: %v", tc.expectedResponce, er),
				)
				return
			}
			if tc.err != "" {
				assert.Equal(
					t,
					tc.err,
					err.Error(),
					fmt.Sprintf("expected response not equal; want %v\n got: %v", tc.err, err.Error()),
				)
				return
			}
			assert.NoError(t, err, "could not remove data from database")
			assert.Equal(
				t,
				tc.expectedResponce,
				er,
				fmt.Sprintf("actual data not equals; want %v\n got: %v", tc.expectedResponce, er),
			)
			_, err = tc.collection.Remove(tc.ctx, tc.enterT.DocumentID)
			assert.NoError(t, err, "could not remove document from database")
		})
	}
}

func TestUpdate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	wrongctx, wrcancel := context.WithTimeout(context.Background(), 5*time.Second)
	wrcancel()
	m, err := Connect(ctx, conURI, dbName, 1, 1)
	require.NoError(t, err, "could not connect to db")
	tt := []struct {
		name             string
		collection       *MariaDB
		enterT           entity.PersonalData
		updateT          entity.PersonalData
		expectedResponce int64
		ctx              context.Context
		err              string
	}{
		{
			name:       "MariaDB -> Update Document in database",
			collection: m,
			enterT: entity.PersonalData{
				DocumentID:  "KGPFOFOGB",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1999,
			},
			updateT: entity.PersonalData{
				DocumentID:  "KGPFOFOGB",
				Name:        "FirstName",
				LastName:    "LastName",
				Phone:       "999999999",
				Email:       "test@test.test",
				YearOfBirth: 1999,
			},
			expectedResponce: 2,
			ctx:              ctx,
			err:              "",
		},
		{
			name:       "MariaDB -> Update Document in database",
			collection: m,
			enterT: entity.PersonalData{
				DocumentID:  "KGPFOFOGB",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1999,
			},
			updateT: entity.PersonalData{
				DocumentID:  "abc",
				Name:        "FirstName",
				LastName:    "LastName",
				Phone:       "999999999",
				Email:       "test@test.test",
				YearOfBirth: 1999,
			},
			expectedResponce: 0,
			ctx:              ctx,
			err:              "",
		},
		{
			name:       "MariaDB -> Faild update document in database",
			collection: m,
			enterT: entity.PersonalData{
				DocumentID:  "KGPFOFOGB",
				Name:        "Name",
				LastName:    "LName",
				Phone:       "1235486",
				Email:       "test@test.test",
				YearOfBirth: 1999,
			},
			updateT:          entity.PersonalData{},
			expectedResponce: 0,
			ctx:              wrongctx,
			err:              "could not update data: context canceled",
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err = tc.collection.Insert(ctx, tc.enterT)
			assert.NoError(t, err, "could not insert data to database")
			er, err := tc.collection.Update(tc.ctx, tc.updateT)
			if tc.expectedResponce != er {
				assert.Equal(
					t,
					tc.expectedResponce,
					er,
					fmt.Sprintf("Update counts not equal; want %v\n got: %v", tc.expectedResponce, er),
				)
				return
			}
			if tc.err != "" {
				assert.Equal(
					t,
					tc.err,
					err.Error(),
					fmt.Sprintf("Errors not equals; want %v\n got: %v", tc.err, err.Error()),
				)
				return
			}
			assert.NoError(t, err, "could not insert data to database")
			_, err = tc.collection.Remove(tc.ctx, tc.enterT.DocumentID)
			assert.NoError(t, err, "could not remove document from database")
		})
	}
}
