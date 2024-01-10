package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"forum/database"
	"forum/urlHandlers"
	"forum/validateData"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database.Engine()
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.TLSConfig = &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		InsecureSkipVerify:       true,
	}
	staticFiles := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFiles))

	mux.HandleFunc("/register", urlHandlers.HandleRegister)
	mux.HandleFunc("/login", urlHandlers.HandleLogin)
	mux.HandleFunc("/logout", urlHandlers.HandleLogout)
	mux.HandleFunc("/", urlHandlers.HandleForum)
	mux.HandleFunc("/post", urlHandlers.HandlePost)
	mux.HandleFunc("/postcontent", urlHandlers.HandlePostContent)
	mux.HandleFunc("/profile", urlHandlers.HandleProfile)
	mux.HandleFunc("/googleAuth", urlHandlers.HandleGoogleAuth)
	mux.HandleFunc("/githubAuth", urlHandlers.HandleGithubAuth)

	fmt.Println("Server hosted at: https://localhost:" + "8080")
	fmt.Println("To Kill Server press Ctrl+C")

	err := server.ListenAndServeTLS("cert.pem", "key.pem")
	validateData.CheckErr(err)
}
