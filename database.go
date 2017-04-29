package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type database struct {
	db *sql.DB
}

func NewDatabase(connection string) (d database, err error) {
	d.db, err = sql.Open("postgres", connection)

	return
}

func (d database) AreValidCredentials(u, p string) (err error) {
	var hashedPassword string

	rows, err := d.db.Query("SELECT password FROM users WHERE username = $1", u)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&hashedPassword)
		if err != nil {
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(p))
		return
	}

	err = rows.Err()
	if err != nil {
		return
	}

	return fmt.Errorf("no such user %q", u)
}
