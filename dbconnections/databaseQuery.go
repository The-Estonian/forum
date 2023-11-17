package dbconnections

import (
	"database/sql"
	"fmt"

	"forum/validateData"
)

func RegisterUser(username, email, password string) (bool, bool) {
	database, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)

	usernameCheck := checkValueFromDB(database, "username", username)
	emailCheck := checkValueFromDB(database, "email", email)
	if !usernameCheck && !emailCheck {
		datastream, err := database.Prepare("INSERT INTO users(username, password, email) VALUES(?, ?, ?)")
		validateData.CheckErr(err)
		datastream.Exec(username, password, email)
		fmt.Println("New user added to the DB")
	}
	database.Close()
	return usernameCheck, emailCheck
}

func LoginUser(username, password string) bool {
	database, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	getPassword := checkPassword(database, username)
	return getPassword == password
}
