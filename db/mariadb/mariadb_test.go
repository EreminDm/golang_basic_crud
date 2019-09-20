package mariadb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	conURI = "test:test@tcp(127.0.0.1:3306)"
	dbName = "person"
)

func TestConnect(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	expected := &MariaDB{}

	tt := []struct {
		name          string
		context       context.Context
		connectionURI string
		dbName        string
		err           string
	}{
		{
			name:          "Success connection",
			context:       ctx,
			connectionURI: conURI,
			dbName:        dbName,
			err:           "",
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			db, err := Connect(tc.context, tc.connectionURI, tc.dbName)
			if tc.err != "" {
				assert.Equal(
					t,
					tc.err,
					err.Error(),
					fmt.Sprintf("expected error is not equals; want: %v, got: %v", tc.err, err.Error()),
				)
				return
			}
			assert.NoError(t, err, "could not connect to mariaDB")
			assert.IsType(t, expected, db, "not equals")

		})
	}

}
