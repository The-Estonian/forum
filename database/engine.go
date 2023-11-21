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
		dbconnections.CreateSessions()
		fmt.Println("Session hash Table Created!")
		dbconnections.CreateCategories()
		fmt.Println("Categories Table Created!")
		dbconnections.CreateLikes()
		fmt.Println("Likes Table Created!")
		dbconnections.CreatePosts()
		fmt.Println("Posts Table Created!")
		fmt.Println("Full Database Regeneration Successfull!")
	}
}
