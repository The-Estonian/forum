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

func checkPassword(database *sql.DB, username string) string {
	var returnString string
	err := database.QueryRow("SELECT password FROM users WHERE username=?", username).Scan(&returnString)
	fmt.Println(err)
	if err != nil {
		fmt.Println("User does not exist")
		return ""
	}
	return returnString
}
