package urlHandlers

import (
	"fmt"
	"forum/dbconnections"
	"forum/helpers"
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
		m.CategoryChoice[0].Selected = "true"
		template.Execute(w, m)
		return
	}
	r.ParseForm()
	if r.Form["Category"] != nil {
		m.AllPosts = helpers.FilterByCat(m, r.Form["Category"][0])
		for i := 0; i < len(m.CategoryChoice); i++ {
			if m.CategoryChoice[i].Category == r.Form["Category"][0] {
				m.CategoryChoice[i].Selected = "true"
			}
		}
		executeCat := template.Execute(w, m)
		if executeCat != nil {
			fmt.Println("Template error: ", executeCat)
		}
	}

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
	m.CategoryChoice[0].Selected = "true"

	executeErr := template.Execute(w, m)
	if executeErr != nil {
		fmt.Println("Template error: ", executeErr)
	}
}
