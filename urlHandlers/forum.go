package urlHandlers

import (
	"fmt"
	"forum/dbconnections"
	"html/template"
	"net/http"
)

func HandleForum(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./templates/forum.html")
	if err != nil {
		http.Error(w, "Template not found!"+err.Error(), http.StatusInternalServerError)
	}
	m := dbconnections.GetMegaDataValues(r, "Forum")
	if r.Method != http.MethodPost {
		template.Execute(w, m)
		return
	}
	r.ParseForm()
	userCurrLike := dbconnections.GetPostLike(m.User.Id, r.Form["postId"][0])
	if r.Form["like"] != nil {
		if userCurrLike == "1" {
			dbconnections.SetPostLikes(m.User.Id, r.Form["postId"][0], "0")
		} else {
			dbconnections.SetPostLikes(m.User.Id, r.Form["postId"][0], "1")
		}
	}
	if r.Form["dislike"] != nil {
		if userCurrLike == "-1" {
			dbconnections.SetPostLikes(m.User.Id, r.Form["postId"][0], "0")
		} else {
			dbconnections.SetPostLikes(m.User.Id, r.Form["postId"][0], "-1")
		}
	}
	m = dbconnections.GetMegaDataValues(r, "Forum")

	executeErr := template.Execute(w, m)
	if executeErr != nil {
		fmt.Println("Template error: ", executeErr)
	}
}
