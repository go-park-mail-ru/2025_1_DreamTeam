package postgres

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	conn           *sql.DB
	SESSION_SECRET string
}

// NewDatabase - конструктор
func NewDatabase(connStr string, SESSION_SECRET string) (*Database, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		err = db.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	_, err = db.Exec(`SET statement_timeout = '3s'; SET lock_timeout = '400ms';`)
	if err != nil {
		err = db.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	return &Database{conn: db, SESSION_SECRET: SESSION_SECRET}, nil
}

// Close - закрытие соединения с базой данных
func (d *Database) Close() {
	err := d.conn.Close()
	if err != nil {
		return
	}
}
