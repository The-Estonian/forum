package urlHandlers

import (
	"fmt"
	"html/template"
	"net/http"

	"forum/cleanData"
	"forum/dbconnections"
	"forum/validateData"
)

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./templates/register.html")
	if err != nil {
		http.Error(w, "Template not found!"+err.Error(), http.StatusInternalServerError)
	}
	if r.Method != http.MethodPost {
		template.Execute(w, nil)
		return
	}

	formDataUsername := r.FormValue("username")
	formDataEmail := r.FormValue("email")
	formDataPassword1 := r.FormValue("password1")
	formDataPassword2 := r.FormValue("password2")

	errorLog := []string{}
	dataValid := true

	passValidFirst, passValidSecond := validateData.ValidatePassword(formDataPassword1, formDataPassword2)
	if !validateData.ValidateName(formDataUsername) {
		errorLog = append(errorLog, "Username should be minimum 3 letters!")
		dataValid = false
	}
	if !validateData.ValidateEmail(formDataEmail) {
		errorLog = append(errorLog, "Please enter a valid email!")
		dataValid = false
	}
	if !passValidFirst {
		errorLog = append(errorLog, "Password should be at least 6 letters long!")
		dataValid = false
	}
	if !passValidSecond {
		errorLog = append(errorLog, "Passwords do not match!")
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
	email := cleanData.CleanEmail(formDataEmail)
	password := formDataPassword1

	userNameOk, userEmailOk := dbconnections.RegisterUser(username, email, password)
	if userNameOk {
		executeErr := template.Execute(w, "Username allready exists!")
		if executeErr != nil {
			fmt.Println("Template error: ", executeErr)
		}
		return
	}
	if userEmailOk {
		executeErr := template.Execute(w, "Email allready exists!")
		if executeErr != nil {
			fmt.Println("Template error: ", executeErr)
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
