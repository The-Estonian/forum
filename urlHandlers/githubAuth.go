package urlHandlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func HandleGithubAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Github time!")
	code := r.URL.Query().Get("code")
	requestBodyMap := map[string]string{
		"client_id":     "c1fe0985284e0a15a4c7",
		"client_secret": "008c5227b7d3b08b2d4358238f398640be335298",
		"code":          code,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)
	tokenRequest, reqerr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON))
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	tokenRequest.Header.Set("Content-Type", "application/json")
	tokenRequest.Header.Set("Accept", "application/json")
	tokenResponse, err := http.DefaultClient.Do(tokenRequest)
	if err != nil {
		log.Panic("Request failed")
	}
	body, _ := io.ReadAll(tokenResponse.Body)
	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}
	var token githubAccessTokenResponse
	json.Unmarshal(body, &token)

	dataRequest, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil)
	if err != nil {
		log.Panic("API Request creation failed")
	}
	authorizationHeaderValue := fmt.Sprintf("token %s", token)
	dataRequest.Header.Set("Authorization", authorizationHeaderValue)

	dataResponse, err := http.DefaultClient.Do(dataRequest)
	if err != nil {
		log.Panic("Request failed")
	}
	responsebody, _ := io.ReadAll(dataResponse.Body)
	fmt.Println(string(responsebody))
}
