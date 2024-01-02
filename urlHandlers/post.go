package urlHandlers

import (
	"forum/dbconnections"
	"html/template"
	"net/http"
)

func HandlePost(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./templates/post.html")
	if err != nil {
		http.Error(w, "Template not found!"+err.Error(), http.StatusInternalServerError)
	}

	m := dbconnections.GetMegaDataValues(r, "Post")
	if r.Method != http.MethodPost {
		template.Execute(w, m)
		return
	}
	
	dbconnections.InsertMessage(r, m.User.Id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
