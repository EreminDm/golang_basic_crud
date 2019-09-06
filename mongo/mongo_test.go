package mongo_test

import (
	"context"
	"fmt"

	"github.com/EreminDm/golang_basic_crud/mongo"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestConnect(t *testing.T) {
	expected := &mongo.Mongodatabase{}
	tt := []struct {
		name          string
		context       context.Context
		connectionURI string
		dbName        string
		err           string
	}{
		{
			name:          "Context TODO",
			context:       context.TODO(),
			connectionURI: "192.168.99.100:27017",
			dbName:        "test",
			err:           "",
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mc, err := mongo.Connect(tc.context, tc.connectionURI, tc.dbName)
			assert.NoError(
				t,
				err,
				fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err),
			)

			if tc.err != "" {
				assert.Equal(
					t,
					tc.err,
					err,
					fmt.Sprintf("expected status Bad Request; got: %v", err),
				)
				return
			}
			assert.IsType(t, expected, mc, fmt.Sprintf("expected type %v, actual %v", expected, mc))
		})
	}
}
