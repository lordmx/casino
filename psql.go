package main

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
)

func NewPsql(
	host string,
	port int,
	user, password, database, names, timezone string,
) (*sql.DB, error) {
	if host == "" {
		host = "localhost"
	}

	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s ",
		user,
		password,
		database,
		host,
	)

	if port != 0 {
		dsn += fmt.Sprintf("port=%d ", port)
	}

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	if names != "" {
		_, err = db.Exec(fmt.Sprintf("SET NAMES '%s'", names))

		if err != nil {
			return nil, err
		}
	}

	if timezone != "" {
		db.Exec(fmt.Sprintf("SET TIME ZONE '%s'", timezone))

		if err != nil {
			return nil, err
		}
	}

	return db, nil
}