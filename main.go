package main

import (
	"fmt"
	"net/http"

	"01.kood.tech/git/jsaar/forum/database"
	"01.kood.tech/git/jsaar/forum/urlHandlers"
	"01.kood.tech/git/jsaar/forum/validateData"
	_ "github.com/mattn/go-sqlite3"
)

var PORT = "8080"

func main() {
	database.Engine()
	http.HandleFunc("/", urlHandlers.HandleIndex)
	http.HandleFunc("/register", urlHandlers.HandleRegister)

	fmt.Println("Server hosted at: http://localhost:" + PORT)
	fmt.Println("To Kill Server press Ctrl+C")
	
	err := http.ListenAndServe(":"+PORT, nil)
	validateData.CheckErr(err)
}
