package services

import (
	sqlite "gosqlite.org"
)

func SetupDbConn() (*sqlite.DB, error) {
	conn, err := sqlite.OpenWAL("animals.db")
	if err != nil {
		return nil, err
	}

	err = MigrateTables(conn)
	return conn, err
}

func MigrateTables(db *sqlite.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS animals (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		name TEXT NOT NULL,
		breed TEXT NOT NULL,
		age INTEGER NOT NULL,
		height INTEGER NOT NULL
	)`)

	return err
}
