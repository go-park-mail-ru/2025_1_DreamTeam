package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Database struct {
	conn *sql.DB
}

// NewDatabase - конструктор
func NewDatabase(connStr string) (*Database, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return &Database{conn: db}, nil
}

// Close - закрытие соединения с базой данных
func (d *Database) Close() {
	d.conn.Close()
}
