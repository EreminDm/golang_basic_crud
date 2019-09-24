package mariadb

import (
	"context"
	"database/sql"
	"time"

	// Register mysql package
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

// MariaDB is connection to mysql database.
type MariaDB struct {
	Person *sql.DB
}

// Connect returns mySQL database connection information,
// do not forget close db connection,
// dataSourcename example: "b_username:db_password@protocol(address:port_num)".
func Connect(ctx context.Context, dataSourceName string, dbname string, maxidleconn, maxopenconn int) (*MariaDB, error) {
	db, err := sql.Open("mysql", dataSourceName+"/"+dbname)
	if err != nil {
		return nil, errors.Wrap(err, "could not open mysql connection")
	}
	err = db.PingContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "could not get stable connection to databases")
	}
	db.SetConnMaxLifetime(time.Minute)
	db.SetMaxIdleConns(maxidleconn)
	db.SetMaxOpenConns(maxopenconn)

	return &MariaDB{
		Person: db,
	}, nil
}
