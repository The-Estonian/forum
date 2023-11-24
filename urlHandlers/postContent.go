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

	// m := structs.MegaData{
	// 	User:        structs.User{Id: "1", Username: "admin", Email: "asd@asd.com", UserAccess: "Bueno!"},
	// 	Post:        post,
	// 	AllComments: dbconnections.GetAllComments(db, postId),
	// }

	m := dbconnections.GetMegaDataValues(r)

	if r.Method != http.MethodPost {
		template.Execute(w, m)
		return
	}
	comment := r.FormValue("createPostComment")

	dbconnections.InsertComment(db, m.Post.Id, m.User.Id, comment)

	executeErr := template.Execute(w, m)
	if executeErr != nil {
		fmt.Println("Template error: ", executeErr)
	}
}
