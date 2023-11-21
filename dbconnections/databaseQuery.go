package dbconnections

import (
	"database/sql"
	"fmt"

	"forum/structs"
	"forum/validateData"
)

func RegisterUser(db *sql.DB, username, email, password string) (bool, bool) {
	usernameCheck := CheckValueFromDB(db, "username", username)
	emailCheck := CheckValueFromDB(db, "email", email)
	if !usernameCheck && !emailCheck {
		_, err := db.Exec("INSERT INTO users(username, password, email) VALUES(?, ?, ?)", username, password, email)
		validateData.CheckErr(err)
		fmt.Println("New user added to the DB")
	}
	return usernameCheck, emailCheck
}

func LoginUser(db *sql.DB, username, password string) bool {
	getPassword := CheckPassword(db, username)
	return getPassword == password
}

func LogoutUser(db *sql.DB, userID string) {
	fmt.Println(userID)
	_, err := db.Exec("DELETE FROM session WHERE user=?", userID)
	validateData.CheckErr(err)
}

func ApplyHash(db *sql.DB, user, hash string) {
	datastream, err := db.Prepare("INSERT OR REPLACE INTO session(user, hash) VALUES(?, ?)")
	validateData.CheckErr(err)
	datastream.Exec(user, hash)
}

func GetID(db *sql.DB, username string) string {
	query := db.QueryRow("SELECT id FROM users WHERE username=?", username).Scan(&username)
	if query != nil {
		fmt.Println("Didn't find username with that name to return ID")
		fmt.Println("Error code: ", query)
	}
	return username
}

func CheckHash(db *sql.DB, hash string) string {
	var user string
	query := db.QueryRow("SELECT user FROM session WHERE hash=?", hash).Scan(&user)
	if query != nil {
		fmt.Println("Didn't find user with that hash!")
		fmt.Println("Error code: ", query)
		return ""
	}
	return user
}

func CheckValueFromDB(db *sql.DB, column string, valueToCheck string) bool {

	newUsername := db.QueryRow("SELECT "+column+" FROM users WHERE "+column+"=?", valueToCheck).Scan(&valueToCheck)
	trigger := false
	if newUsername == nil {
		fmt.Println("Username already exists!")
		trigger = true
	}
	return trigger
}

func CheckPassword(db *sql.DB, username string) string {
	var returnString string
	err := db.QueryRow("SELECT password FROM users WHERE username=?", username).Scan(&returnString)
	if err != nil {
		fmt.Println("User does not exist")
		return ""
	}
	return returnString
}

func GetAllPosts(db *sql.DB) []structs.Post {
	var allPosts []structs.Post
	posts, _ := db.Query("SELECT * FROM posts")
	for posts.Next() {
		var post structs.Post
		if err := posts.Scan(&post.Id, &post.Category, &post.User, &post.Post, &post.Created); err != nil {
			fmt.Println(err)
			return allPosts
		}
		allPosts = append(allPosts, post)
	}
	return allPosts
}
