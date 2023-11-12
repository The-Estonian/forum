package dbconnections

import (
	"database/sql"
	"fmt"
)

func checkValueFromDB(database *sql.DB, column string, valueToCheck string) bool {

	newUsername := database.QueryRow("SELECT "+column+" FROM users WHERE "+column+"=?", valueToCheck).Scan(&valueToCheck)
	trigger := false
	if newUsername == nil {
		fmt.Println("Username already exists!")
		trigger = true
	}
	return trigger
}
