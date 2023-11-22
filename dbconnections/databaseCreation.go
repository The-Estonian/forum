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
	db, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	db.Exec("CREATE TABLE `users` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `username` VARCHAR(255) NOT NULL, `password` VARCHAR(255) NOT NULL, `email` VARCHAR(255) NOT NULL)")
	db.Close()
}

func CreateSessions() {
	database, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	database.Exec("CREATE TABLE `session` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `user` INTEGER UNIQUE REFERENCES users(id), `hash` VARCHAR(255) NOT NULL)")
	database.Close()
}

func CreateCategory() {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	db.Exec("CREATE TABLE `category` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `category` VARCHAR(255) NOT NULL)")

	for _, v := range []string{"Potato", "Carrot", "Tomatoe", "Apple", "Orange"} {
		_, err := db.Exec("INSERT INTO category (category) VALUES (?)", v)
		if err != nil {
			fmt.Println(err)
		}
	}

	db.Close()
}

func CreatePostCategoryList() {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	db.Exec("CREATE TABLE `post_category_list` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `post_category` INTEGER NOT NULL REFERENCES category(id), `post_id` INTEGER NOT NULL REFERENCES posts(id))")
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
	database.Exec("CREATE TABLE `posts` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `title` VARCHAR(255) NOT NULL, `user` INTEGER NOT NULL REFERENCES users(id), `post` VARCHAR(255), `created` NOT NULL DEFAULT CURRENT_TIMESTAMP)")
	database.Close()
}
