package dbconnections

import (
	"database/sql"
	"fmt"
	"os"

	"01.kood.tech/git/jsaar/forum/checkErrors"
)

func CreateDB() {
	if err := os.MkdirAll("./database", os.ModeSticky|os.ModePerm); err != nil {
		fmt.Println(err)
	}
	os.Create("./database/forum.db")
}

func CreateTables() {
	database, err := sql.Open("sqlite3", "./database/forum.db")
	checkErrors.CheckErr(err)
	database.Exec("CREATE TABLE `users` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `username` VARCHAR(255) NOT NULL, `password` VARCHAR(255) NOT NULL, `email` VARCHAR(255) NOT NULL)")
	database.Close()
}

func RegisterUser(username, password, email string) {
	database, err := sql.Open("sqlite3", "./database/forum.db")
	checkErrors.CheckErr(err)
	datastream, err := database.Prepare("INSERT INTO users(username, password, email) VALUES(?, ?, ?)")
	datastream.Exec(username, password, email)
	database.Close()
}

// func sendToDB() {
// 	database, err := sql.Open("sqlite3", "./database/forum.db")
// 	checkErrors.CheckErr(err)
// 	datastream, err := database.Prepare("INSERT INTO users")
// 	database.Query()
// }
