package mariadb

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

// MariaDB is connection to mysql database.
type MariaDB struct {
	Person *sql.DB
}

// Connect returns mySQL database connection information,
// do not forget close db connection,
// dataSourcename example: "b_username:db_password@protocol(address:port_num)".
func Connect(ctx context.Context, dataSourceName string, dbname string) (*MariaDB, error) {

	db, err := sql.Open("mysql", dataSourceName+"/"+dbname)
	if err != nil {
		return nil, errors.Wrap(err, "could not open mysql connection")
	}
	err = db.PingContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "could not get stable connection to databases")
	}
	db.SetConnMaxLifetime(time.Minute)
	db.SetMaxIdleConns(3)
	db.SetMaxOpenConns(30)

	return &MariaDB{
		Person: db,
	}, nil
}