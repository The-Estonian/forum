package dbconnections

import (
	"database/sql"
	"fmt"
	"os"

	"forum/validateData"
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
