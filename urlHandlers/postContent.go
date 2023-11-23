package urlHandlers

import (
	"database/sql"
	"fmt"
	"forum/dbconnections"
	"forum/validateData"
	"html/template"
	"net/http"
)

func HandlePostContent(w http.ResponseWriter, r *http.Request) {

	template, err := template.ParseFiles("./templates/postContent.html")
	if err != nil {
		http.Error(w, "Template not found!"+err.Error(), http.StatusInternalServerError)
	}

	db, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	defer db.Close()

	cookie, err := r.Cookie("UserCookie")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !dbconnections.HashInDatabase(db, cookie.Value) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postId := r.URL.Query().Get("id")
	post := dbconnections.GetOnePost(db, postId)

	if r.Method != http.MethodPost {
		template.Execute(w, post)
		return
	}
	commentatorId := dbconnections.CheckHash(db, cookie.Value)
	comment := r.FormValue("createPostComment")

	dbconnections.InsertComment(db, postId, commentatorId, comment)

	executeErr := template.Execute(w, post)
	if executeErr != nil {
		fmt.Println("Template error: ", executeErr)
	}
}
