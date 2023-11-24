package dbconnections

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"

	"forum/structs"
	"forum/validateData"
)

// Open and Close database connection. Returns database connection.
func DbConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	return db
}

// Register new user, add user values to users database. Returns true if user/email is database.
func RegisterUser(db *sql.DB, username, email, password string) (bool, bool) {
	usernameCheck := CheckValueFromDB(db, "username", username)
	emailCheck := CheckValueFromDB(db, "email", email)
	if !usernameCheck && !emailCheck {
		_, err := db.Exec("INSERT INTO users(username, password, email) VALUES(?, ?, ?)", username, password, email)
		validateData.CheckErr(err)
		fmt.Println("New user added to the DB")
		SetAccessRight(GetID(username), "2")
		fmt.Println("Access granted to user", GetID(username))
	}
	return usernameCheck, emailCheck
}

// Returns True if user inserted credentials are in database.
func LoginUser(username, password string) bool {
	getPassword := CheckPassword(username)
	return getPassword == password
}

// Deletes Cookie from session database
func LogoutUser(db *sql.DB, userID string) {
	fmt.Println(userID)
	_, err := db.Exec("DELETE FROM session WHERE user=?", userID)
	validateData.CheckErr(err)
}

// Applies Cookie in session database
func ApplyHash(user, hash string) {
	db := DbConnection()
	datastream, err := db.Prepare("INSERT OR REPLACE INTO session(user, hash) VALUES(?, ?)")
	defer db.Close()
	validateData.CheckErr(err)
	datastream.Exec(user, hash)
}

// Returns ID if username exists in the users database
func GetID(username string) string {
	db := DbConnection()
	query := db.QueryRow("SELECT id FROM users WHERE username=?", username).Scan(&username)
	defer db.Close()
	if query != nil {
		fmt.Println("Didn't find username with that name to return ID")
		fmt.Println("Error code: ", query)
	}
	return username
}

// Returns All user info by user hash
func GetUserInfo(id string) structs.User {
	db := DbConnection()
	var userInfo structs.User
	var dump string
	query := db.QueryRow("SELECT * FROM users WHERE id=?", id).Scan(&userInfo.Id, &userInfo.Username, &dump, &userInfo.Email)
	dump = ""
	defer db.Close()
	if query != nil {
		fmt.Println("Didn't find userId with that id")
		fmt.Println("Error code: ", query)
	}
	return userInfo
}

// Returns the users access rights
func GetAccessRight(id string) structs.AccessRights {
	db := DbConnection()
	var userAccess structs.AccessRights
	query := db.QueryRow("SELECT user_access FROM user_access WHERE user=?", id).Scan(&userAccess.AccessRight)
	defer db.Close()
	if query != nil {
		fmt.Println("Didn't find user with that id")
		fmt.Println("Error code: ", query)
	}
	return userAccess
}

// Sets the users access rights
func SetAccessRight(user string, access string) {
	db := DbConnection()
	_, err := db.Exec("INSERT INTO user_access (user, user_access) VALUES(?, ?)", user, access)
	if err != nil {
		fmt.Println("SetAccessRight")
		fmt.Println("Error code: ", err)
	}
	defer db.Close()
}

// Returns UserID from session database
func CheckHash(hash string) string {
	db := DbConnection()
	var user string
	query := db.QueryRow("SELECT user FROM session WHERE hash=?", hash).Scan(&user)
	defer db.Close()
	if query != nil {
		fmt.Println("CheckHash")
		fmt.Println("Error code: ", query)
		return "1"
	}
	return user
}

// Returns True if hash in database
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

// Return True if value is in users database column
func CheckValueFromDB(db *sql.DB, column string, valueToCheck string) bool {
	newUsername := db.QueryRow("SELECT "+column+" FROM users WHERE "+column+"=?", valueToCheck).Scan(&valueToCheck)
	trigger := false
	if newUsername == nil {
		fmt.Println("Username already exists!")
		trigger = true
	}
	return trigger
}

// Returns password from users based on username. Hopefully encrypted.
func CheckPassword(username string) string {
	db := DbConnection()
	var returnString string
	err := db.QueryRow("SELECT password FROM users WHERE username=?", username).Scan(&returnString)
	defer db.Close()
	if err != nil {
		fmt.Println("User does not exist")
		return ""
	}
	return returnString
}

// Returns all rows in an array of structs from posts database
func GetAllPosts() []structs.Post {
	var allPosts []structs.Post
	db := DbConnection()
	posts, _ := db.Query("SELECT * FROM posts")
	defer db.Close()
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

// Returns a struct that contains data from one row in post database
func GetOnePost(data string) structs.Post {
	db := DbConnection()
	posts := db.QueryRow("SELECT * FROM posts WHERE id=?", data)
	defer db.Close()
	var post structs.Post
	if err := posts.Scan(&post.Id, &post.Title, &post.User, &post.Post, &post.Created); err != nil {
		fmt.Println(err)
	}
	return post
}

// Inserts into posts and post_category_list user inserted data
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

// Inserts comments database user inserted comment and commentator userID
func InsertComment(db *sql.DB, postId string, commentatorId string, comment string) {
	_, err := db.Exec("INSERT INTO comments (post_id, user, comment) VALUES (?, ?, ?)", postId, commentatorId, comment)
	if err != nil {
		fmt.Println(err)
	}
}

// Returns all comments by post_id
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

func GetMegaDataValues(r *http.Request) structs.MegaData {
	var userId string
	cookie, err := r.Cookie("UserCookie")
	if err != nil {
		userId = "1"
	} else {
		userId = CheckHash(cookie.Value)
	}

	if len(r.URL.Query().Get("id")) > 0 {
		fmt.Println("URL QUERY DATA LEN: ", len(r.URL.Query().Get("id")))
		fmt.Println(r.URL.Query().Get("id"))
		postId := r.URL.Query().Get("id")
		m := structs.MegaData{
			User:   GetUserInfo(userId),
			Post:   GetOnePost(postId),
			Access: GetAccessRight(userId),
		}
		return m
	}

	m := structs.MegaData{
		User:     GetUserInfo(userId),
		AllPosts: GetAllPosts(),
		Access:   GetAccessRight(userId),
	}

	return m
}
