package urlHandlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"forum/dbconnections"
	"forum/validateData"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./templates/logout.html")
	if err != nil {
		http.Error(w, "Template not found!"+err.Error(), http.StatusInternalServerError)
	}
	db, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	defer db.Close()
	cookie, err := r.Cookie("UserCookie")
	if err == nil {
		hash := dbconnections.CheckHash(db, cookie.Value)
		fmt.Println(hash)
		dbconnections.LogoutUser(db, hash)
	}
	exp := time.Now().Add(1 * time.Millisecond)
	cookie = &http.Cookie{
		Name:     "UserCookie",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  exp,
	}
	http.SetCookie(w, cookie)

	executeErr := template.Execute(w, nil)
	if executeErr != nil {
		fmt.Println("Template error: ", executeErr)
	}
}
