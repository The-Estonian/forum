package urlHandlers

import (
	"fmt"
	"forum/cleanData"
	"forum/dbconnections"
	"forum/validateData"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./templates/login.html")
	if err != nil {
		http.Error(w, "Template not found!"+err.Error(), http.StatusInternalServerError)
	}

	m := dbconnections.GetMegaDataValues(r, "Login")

	if r.Method != http.MethodPost {
		template.Execute(w, m)
		return
	}

	r.ParseMultipartForm(0)
	if r.FormValue("loginType") == "google" {
		authURL := fmt.Sprintf(
			"https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email&state=%s",
			url.QueryEscape("113897564750-m48e0sfio9uiei7k9v6aipli8526q97t.apps.googleusercontent.com"),
			url.QueryEscape("https://localhost:8080/googleAuth"),
			url.QueryEscape("ForumAuthentication"))
		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
		return
	}
	if r.FormValue("loginType") == "github" {
		clientID := "c1fe0985284e0a15a4c7"
		redirectURI := "https://localhost:8080/githubAuth"
		authURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", url.QueryEscape(clientID), url.QueryEscape(redirectURI))
		fmt.Println(authURL)
		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
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
		m.Errors = errorLog
		executeErr := template.Execute(w, m)
		if executeErr != nil {
			fmt.Println("Template error: ", executeErr)
		}
		return
	}

	username := cleanData.CleanName(formDataUsername)

	if dbconnections.LoginUser(username, formDataPassword) {
		id := uuid.New()
		exp := time.Now().Add(10 * time.Minute)
		cookie := &http.Cookie{
			Name:     "UserCookie",
			Value:    id.String(),
			Path:     "/",
			HttpOnly: true,
			Expires:  exp}
		http.SetCookie(w, cookie)
		dbconnections.ApplyHash(dbconnections.GetID(username), id.String())

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	m.Errors = append(m.Errors, "Username or Password incorrect")
	executeErr := template.Execute(w, m)
	if executeErr != nil {
		fmt.Println("Template error: ", executeErr)
	}
}
