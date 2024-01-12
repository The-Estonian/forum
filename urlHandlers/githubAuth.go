package urlHandlers

import (
	"encoding/json"
	"fmt"
	"forum/cleanData"
	"forum/dbconnections"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

func HandleGithubAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Github time!")
	code := r.URL.Query().Get("code")
	clientID := "c1fe0985284e0a15a4c7"
	githubTokenEndpoint := "https://github.com/login/oauth/access_token"
	clientSecret := "dd70f816d1d6be3c7386dd69591fe548455ac193"
	redirectURL := "https://localhost:8080/githubAuth"
	tokenURL := fmt.Sprintf("%s?code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code",
		githubTokenEndpoint, url.QueryEscape(code), url.QueryEscape(clientID), url.QueryEscape(clientSecret), url.QueryEscape(redirectURL))

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
	params := strings.Split(string(body), "&")
	paramsMap := make(map[string]string)
	for _, param := range params {
		keyValue := strings.Split(param, "=")
		if len(keyValue) == 2 {
			paramsMap[keyValue[0]] = keyValue[1]
		}
	}
	accessToken := paramsMap["access_token"]
	emailRequest, err := http.NewRequest(
		"GET",
		"https://api.github.com/user/emails",
		nil)
	if err != nil {
		fmt.Println("Error setting new http request", err)
	}
	emailRequest.Header.Set("Authorization", "Bearer "+accessToken)
	emailResponse, err := http.DefaultClient.Do(emailRequest)
	if err != nil {
		fmt.Println("Error setting new http request", err)
	}
	defer emailResponse.Body.Close()
	emailBody, err := io.ReadAll(emailResponse.Body)
	if err != nil {
		fmt.Println("Error setting new http request", err)
	}
	var emails []map[string]interface{}
	if err := json.Unmarshal(emailBody, &emails); err != nil {
		fmt.Println("Error setting new http request", err)
	}
	emailCheck := cleanData.CleanEmail((emails[0]["email"]).(string))
	if len(emailCheck) > 0 {
		username := dbconnections.GetUsername(emailCheck)
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
