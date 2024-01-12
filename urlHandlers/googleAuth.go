package urlHandlers

import (
	"encoding/json"
	"fmt"
	"forum/cleanData"
	"forum/dbconnections"
	"forum/helpers"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

func HandleGoogleAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Google time!")
	clientID := "113897564750-m48e0sfio9uiei7k9v6aipli8526q97t.apps.googleusercontent.com"
	googleTokenEndpoint := "https://accounts.google.com/o/oauth2/token"
	clientSecret := "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"
	redirectURL := "https://localhost:8080/googleAuth"

	code := r.FormValue("code")
	tokenURL := fmt.Sprintf("%s?code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code",
		googleTokenEndpoint, url.QueryEscape(code), url.QueryEscape(clientID), url.QueryEscape(clientSecret), url.QueryEscape(redirectURL))

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(""))
	if err != nil {
		fmt.Println("Token request failed:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read token response:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	var responseBody map[string]interface{}
	if err := json.Unmarshal([]byte(string(body)), &responseBody); err != nil {
		fmt.Println("Failed to parse JSON response:", err)
		// Handle the error as needed
		return
	}

	userEmail, err := helpers.ExtractEmailFromIDToken(string(responseBody["id_token"].(string)))
	if err != nil {
		fmt.Println("Email not received")
	}
	userEmail = cleanData.CleanEmail(userEmail)
	if len(userEmail) > 0 {
		username := dbconnections.GetUsername(userEmail)
		// login user
		if len(username) > 0 {
			id := uuid.New()
			exp := time.Now().Add(10 * time.Minute)
			cookie := &http.Cookie{
				Name:     "UserCookie",
				Value:    id.String(),
				Path:     "/",
				HttpOnly: true,
				Expires:  exp}
			http.SetCookie(w, cookie)
			dbconnections.ApplyHash(dbconnections.GetID(username), id.String())
		} else {
			http.Redirect(w, r, "/register?notRegistered=true", http.StatusSeeOther)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

}
