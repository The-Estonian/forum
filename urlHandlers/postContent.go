package urlHandlers

import (
	"fmt"
	"forum/dbconnections"
	"html/template"
	"net/http"
)

func HandlePostContent(w http.ResponseWriter, r *http.Request) {

	template, err := template.ParseFiles("./templates/postContent.html")
	if err != nil {
		http.Error(w, "Template not found!"+err.Error(), http.StatusInternalServerError)
	}

	m := dbconnections.GetMegaDataValues(r, "PostContent")
	if r.Method != http.MethodPost {
		template.Execute(w, m)
		return
	}

	postId := r.URL.Query().Get("PostId")
	commentatorId := m.User.Id
	comment := r.FormValue("createPostComment")

	dbconnections.InsertComment(postId, commentatorId, comment)

	m = dbconnections.GetMegaDataValues(r, "PostContent")
	executeErr := template.Execute(w, m)
	if executeErr != nil {
		fmt.Println("Template error: ", executeErr)
	}
}
