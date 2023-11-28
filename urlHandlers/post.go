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

	// var allCat []structs.Category
	// rows, _ := db.Query("SELECT * FROM category")
	// for rows.Next() {
	// 	var cat structs.Category
	// 	if err := rows.Scan(&cat.Id, &cat.Category); err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	allCat = append(allCat, cat)
	// }

	if r.Method != http.MethodPost {
		template.Execute(w, m)
		return
	}
	r.ParseForm()

	db := dbconnections.DbConnection()
	dbconnections.InsertMessage(db, r.Form, m.User.Id)
	defer db.Close()
	template.Execute(w, m)
}
