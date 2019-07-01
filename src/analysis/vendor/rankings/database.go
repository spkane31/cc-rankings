package rankings

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "postgres"
	dbname = "rankings_dev"
)

func CreateConnectionString() string {
	str := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	return str
}

func ConnectToPSQL() (*sql.DB, error) {
	db, err := sql.Open("postgres", CreateConnectionString())
	return db, err
}