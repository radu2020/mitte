package requests

import (
	"../models"
	"encoding/json"
	"net/http"
)

type CreateUserResponse struct {
	User models.User `json:"user"`
}

type ProfilesResponse struct {
	PotentialMatches []models.User `json:"potential_matches"`
}

type SwipeResponse struct {
	Match bool `json:"match"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ApiKey string `json:"api_key"`
}

type ApiKeyBody struct {
	ApiKey string `json:"api_key"`
}

func HandleRequestNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	resp := make(map[string]string)
	resp["message"] = "Resource Not Found"
	json.NewEncoder(w).Encode(resp)
	return
}

func HandleRequestForbidden(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	resp := make(map[string]string)
	resp["message"] = "Access denied"
	json.NewEncoder(w).Encode(resp)
	return
}
