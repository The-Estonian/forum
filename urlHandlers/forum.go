package urlHandlers

import (
	"database/sql"
	"fmt"
	"forum/dbconnections"
	"forum/validateData"
	"html/template"
	"net/http"
)

func HandleForum(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./templates/forum.html")
	if err != nil {
		http.Error(w, "Template not found!"+err.Error(), http.StatusInternalServerError)
	}

	cookie, err := r.Cookie("UserCookie")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	db, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	defer db.Close()
	data := "Currently logged in user ID: " + dbconnections.CheckHash(db, cookie.Value)

	if r.Method != http.MethodPost {
		template.Execute(w, nil)
		return
	}

	inputData := r.FormValue("data")
	fmt.Println("INPUT DATA: ", inputData)

	for _, v := range dbconnections.GetAllPosts(db) {
		fmt.Println(v)
	}

	executeErr := template.Execute(w, data)
	if executeErr != nil {
		fmt.Println("Template error: ", executeErr)
	}
}
