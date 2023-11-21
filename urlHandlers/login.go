package urlHandlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"forum/cleanData"
	"forum/dbconnections"
	"forum/validateData"

	"github.com/google/uuid"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./templates/login.html")
	if err != nil {
		http.Error(w, "Template not found!"+err.Error(), http.StatusInternalServerError)
	}
	if r.Method != http.MethodPost {
		template.Execute(w, nil)
		return
	}

	formDataUsername := r.FormValue("username")
	formDataPassword := r.FormValue("password")

	errorLog := []string{}
	dataValid := true
	validatePassword, _ := validateData.ValidatePassword(formDataPassword, formDataPassword)
	if !validateData.ValidateName(formDataUsername) {
		errorLog = append(errorLog, "Username should be minimum 3 letters!")
		dataValid = false
	}
	if !validatePassword {
		errorLog = append(errorLog, "Password should be at least 6 letters long!")
		dataValid = false
	}
	if !dataValid {
		executeErr := template.Execute(w, errorLog)
		if executeErr != nil {
			fmt.Println("Template error: ", executeErr)
		}
		return
	}

	username := cleanData.CleanName(formDataUsername)
	db, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	defer db.Close()
	if dbconnections.LoginUser(db, username, formDataPassword) {
		id := uuid.New()
		exp := time.Now().Add(10 * time.Minute)
		cookie := &http.Cookie{
			Name:     "UserCookie",
			Value:    id.String(),
			Path:     "/",
			HttpOnly: true,
			Expires:  exp}
		http.SetCookie(w, cookie)
		dbconnections.ApplyHash(db, dbconnections.GetID(db, username), id.String())

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	executeErr := template.Execute(w, "Username or Password incorrect")
	if executeErr != nil {
		fmt.Println("Template error: ", executeErr)
	}
}
