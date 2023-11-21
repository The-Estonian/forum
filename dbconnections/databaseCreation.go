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

func CreateUsers() {
	database, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	database.Exec("CREATE TABLE `users` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `username` VARCHAR(255) NOT NULL, `password` VARCHAR(255) NOT NULL, `email` VARCHAR(255) NOT NULL)")
	database.Close()
}

func CreateSessions() {
	database, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	database.Exec("CREATE TABLE `session` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `user` INTEGER UNIQUE REFERENCES users(id), `hash` VARCHAR(255) NOT NULL)")
	database.Close()
}

func CreateCategories() {
	database, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	database.Exec("CREATE TABLE `categories` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `category` VARCHAR(255) NOT NULL)")
	database.Close()
}

func CreateLikes() {
	database, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	database.Exec("CREATE TABLE `likes` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `post` INTEGER NOT NULL REFERENCES posts(id), `user` INTEGER NOT NULL REFERENCES users(id))")
	database.Close()
}

func CreatePosts() {
	database, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	database.Exec("CREATE TABLE `posts` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `category` INTEGER NOT NULL REFERENCES categories(id), `user` INTEGER NOT NULL REFERENCES users(id), `post` VARCHAR(255), `created` NOT NULL DEFAULT CURRENT_TIMESTAMP)")
	database.Close()
}
