package urlHandlers

import (
	"fmt"
	"forum/dbconnections"
	"html/template"
	"net/http"
)

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./templates/profile.html")
	if err != nil {
		http.Error(w, "Template not found!"+err.Error(), http.StatusInternalServerError)
	}

	m := dbconnections.GetMegaDataValues(r, "Profile")

	if r.Method != http.MethodPost {
		template.Execute(w, m)
		return
	}

	executeErr := template.Execute(w, m)
	if executeErr != nil {
		fmt.Println("Template error: ", executeErr)
	}
}
