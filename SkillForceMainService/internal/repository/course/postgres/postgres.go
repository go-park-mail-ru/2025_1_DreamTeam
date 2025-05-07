package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Database struct {
	conn           *sql.DB
	SESSION_SECRET string
}

func NewDatabase(connStr string, SESSION_SECRET string) (*Database, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return &Database{conn: db, SESSION_SECRET: SESSION_SECRET}, nil
}

func (d *Database) Close() {
	d.conn.Close()
}
