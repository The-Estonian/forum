package database

import (
	"fmt"
	"os"

	"forum/dbconnections"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func Engine() {
	if !fileExists("./database/forum.db") {
		fmt.Println("Did not find the Database! Starting regeneration!")
		dbconnections.CreateDB()
		fmt.Println("Database Created!")
		dbconnections.CreateUsers()
		fmt.Println("User Table Created!")
		fmt.Println("Full Database Generation Successfull!")
	}
}
