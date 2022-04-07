package test

import (
	"database/sql"
	"fmt"
)

type Table struct {
	Name string `db:"TABLE_NAME"`
}

// Truncate will truncate the all tables given database name
func Truncate(db *sql.DB, databaseName string) error {
	_, err := db.Exec("SET FOREIGN_KEY_CHECKS=0")

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	result, err := db.Query("SELECT TABLE_NAME FROM information_schema.tables WHERE TABLE_SCHEMA = ?", databaseName)

	if err != nil {
		return err
	}

	truncateAll(result, db)

	return nil
}

func truncateAll(rows *sql.Rows, db *sql.DB) {
	for rows.Next() {
		var t = &Table{}

		err := rows.Scan(&t.Name)

		if err != nil {
			fmt.Println(err.Error())

			break
		}

		_, err = db.Exec(fmt.Sprintf("TRUNCATE %v", t.Name))

		if err != nil {
			fmt.Println(err.Error())

			break
		}
	}
}
