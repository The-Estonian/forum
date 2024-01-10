package helpers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"forum/structs"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func contains(elems []structs.Category, v string) bool {
	for _, s := range elems {
		if v == s.Category {
			return true
		}
	}
	return false
}

func FilterByCat(m structs.MegaData, value string) []structs.Post {
	var newPosts []structs.Post
	for i := 0; i < len(m.AllPosts); i++ {
		if contains(m.AllPosts[i].Categories, value) {
			newPosts = append(newPosts, m.AllPosts[i])
		}
	}
	return newPosts
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ExtractEmailFromIDToken(idToken string) (string, error) {
	// Split the ID Token into its three parts (header, payload, signature)
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid ID Token format")
	}

	// Decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("failed to decode ID Token payload: %v", err)
	}

	// Parse the JSON payload
	var payloadData map[string]interface{}
	if err := json.Unmarshal(payload, &payloadData); err != nil {
		return "", fmt.Errorf("failed to parse ID Token payload: %v", err)
	}

	// Extract the "email" claim from the payload
	email, ok := payloadData["email"].(string)
	if !ok {
		return "", fmt.Errorf("email not found in ID Token payload")
	}

	return email, nil
}
