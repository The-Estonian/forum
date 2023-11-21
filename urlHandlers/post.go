package urlHandlers

import (
	"database/sql"
	"fmt"
	"forum/dbconnections"
	"forum/structs"
	"forum/validateData"
	"html/template"
	"net/http"
)

func HandlePost(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./templates/post.html")
	if err != nil {
		http.Error(w, "Template not found!"+err.Error(), http.StatusInternalServerError)
	}

	db, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	defer db.Close()

	cookie, err := r.Cookie("UserCookie")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// if len(dbconnections.CheckHash(db, cookie.Value)) > 0 {
	// 	http.Redirect(w, r, "/login", http.StatusSeeOther)
	// 	return
	// }

	if r.Method != http.MethodPost {
		template.Execute(w, nil)
		return
	}
	var allCat []structs.Category
	rows, err := db.Query("SELECT * FROM categories")
	fmt.Println(err)
	for rows.Next() {
		var cat structs.Category
		if err := rows.Scan(&cat.Id, &cat.Category); err != nil {
			fmt.Println(err)
		}
		allCat = append(allCat, cat)
	}
	for _, v := range allCat {
		fmt.Println(v.Category)

	}

	inputData := r.FormValue("post") //kontrollida üle mis nimi

	dbconnections.InsertPost(db, inputData, dbconnections.CheckHash(db, cookie.Value))

	template.Execute(w, allCat)
}
