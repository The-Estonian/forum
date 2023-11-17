package main

import (
	"fmt"
	"net/http"

	"forum/database"
	"forum/urlHandlers"
	"forum/validateData"
	_ "github.com/mattn/go-sqlite3"
)

var PORT = "8080"

func main() {
	database.Engine()

	staticFiles := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticFiles))

	http.HandleFunc("/", urlHandlers.HandleIndex)
	http.HandleFunc("/register", urlHandlers.HandleRegister)
	http.HandleFunc("/login", urlHandlers.HandleLogin)

	fmt.Println("Server hosted at: http://localhost:" + PORT)
	fmt.Println("To Kill Server press Ctrl+C")

	err := http.ListenAndServe(":"+PORT, nil)
	validateData.CheckErr(err)
}
