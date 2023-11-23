package urlHandlers

import (
	"database/sql"
	"fmt"
	"forum/dbconnections"
	"forum/structs"
	"forum/validateData"
	"html/template"
	"net/http"
	"strings"
)

func HandleForum(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./templates/forum.html")
	if err != nil {
		http.Error(w, "Template not found!"+err.Error(), http.StatusInternalServerError)
	}

	if r.URL.Path != "/" && r.URL.Path != "/register" && r.URL.Path != "/locations" && !strings.HasPrefix(r.URL.Path, "/locations") {
		http.Error(w, "Bad Request: 404", http.StatusNotFound)
		return
	}

	db, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	defer db.Close()
	m := structs.MegaData{
		User:     structs.User{Id: "1", Username: "admin", Email: "asd@asd.com", UserAccess: "Bueno!"},
		AllPosts: dbconnections.GetAllPosts(db),
	}
	executeErr := template.Execute(w, m)
	if executeErr != nil {
		fmt.Println("Template error: ", executeErr)
	}
}
