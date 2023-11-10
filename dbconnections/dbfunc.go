package dbconnections

import (
	"database/sql"
	"fmt"
	"os"

	"01.kood.tech/git/jsaar/forum/validateData"
)

func CreateDB() {
	if err := os.MkdirAll("./database", os.ModeSticky|os.ModePerm); err != nil {
		fmt.Println(err)
	}
	os.Create("./database/forum.db")
}

func CreateTables() {
	database, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	database.Exec("CREATE TABLE `users` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `username` VARCHAR(255) NOT NULL, `password` VARCHAR(255) NOT NULL, `email` VARCHAR(255) NOT NULL)")
	database.Close()
}

func RegisterUser(username, email, password string) (bool, bool) {
	database, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	fmt.Println("Entering: ", email)
	nameCheck := false
	emailCheck := false

	newUsername := database.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&username)

	if newUsername == nil {
		fmt.Println("Username already exists!")
		nameCheck = true
	}
	newEmail := database.QueryRow("SELECT email FROM users WHERE email=?", email).Scan(&email)

	if newEmail == nil {
		fmt.Println("Email already exists!")
		emailCheck = true
	}

	if !nameCheck && !emailCheck {
		datastream, err := database.Prepare("INSERT INTO users(username, password, email) VALUES(?, ?, ?)")
		validateData.CheckErr(err)
		datastream.Exec(username, password, email)
		fmt.Println("New user added to the DB")
	} else {
		if nameCheck {
			nameCheck = true
			fmt.Println("Name exists!")
		}
		if emailCheck {
			emailCheck = true
			fmt.Println("Email exists!")
		}
	}
	database.Close()
	return nameCheck, emailCheck
}
