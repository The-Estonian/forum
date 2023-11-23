package dbconnections

import (
	"database/sql"
	"fmt"
	"net/url"

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

// Testing hash
func HashInDatabase(db *sql.DB, hash string) bool {
	var user string
	query := db.QueryRow("SELECT user FROM session WHERE hash=?", hash).Scan(&user)
	if query != nil {
		fmt.Println("HashInDatabase: didn't find user with that hash!")
		fmt.Println("Error code: ", query)
		return false
	}
	return true
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
		if err := posts.Scan(&post.Id, &post.Title, &post.User, &post.Post, &post.Created); err != nil {
			fmt.Println(err)
			return allPosts
		}
		allPosts = append(allPosts, post)
	}
	return allPosts
}

func GetOnePost(db *sql.DB, data string) structs.Post {
	posts := db.QueryRow("SELECT * FROM posts WHERE id=?", data)

	var post structs.Post
	if err := posts.Scan(&post.Id, &post.Title, &post.User, &post.Post, &post.Created); err != nil {
		fmt.Println(err)
	}
	return post
}

func InsertMessage(db *sql.DB, userForm url.Values, userId string) {

	var inputTitle string
	var inputMessage string
	var catArray []string

	for key, value := range userForm {
		if key == "title" {
			inputTitle = value[0]
		} else if key == "message" {
			inputMessage = value[0]
		} else {
			catArray = append(catArray, key)
		}
	}

	var data int
	err := db.QueryRow("INSERT INTO posts (title, user, post) VALUES (?, ?, ?) RETURNING id", inputTitle, userId, inputMessage).Scan(&data)
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range catArray {
		_, err = db.Exec("INSERT INTO post_category_list (post_category, post_id) VALUES (?, ?)", v, data)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func InsertComment(db *sql.DB, postId string, commentatorId string, comment string) {
	_, err := db.Exec("INSERT INTO comments (post_id, user, comment) VALUES (?, ?, ?)", postId, commentatorId, comment)
	if err != nil {
		fmt.Println(err)
	}
}

func GetAllComments(db *sql.DB, data string) []structs.Comment {
	var allComments []structs.Comment

	allCommentsFromData, _ := db.Query("SELECT * FROM comments WHERE post_id=?", data)
	for allCommentsFromData.Next() {
		var comments structs.Comment
		if err := allCommentsFromData.Scan(&comments.Id, &comments.PostId, &comments.UserId, &comments.Comment, &comments.Created); err != nil {
			fmt.Println(err)
		}
		allComments = append(allComments, comments)
	}

	return allComments

}
