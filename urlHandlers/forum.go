package urlHandlers

import (
	"fmt"
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

	executeErr := template.Execute(w, cookie.Expires)
	if executeErr != nil {
		fmt.Println("Template error: ", executeErr)
	}
}
