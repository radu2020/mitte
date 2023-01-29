package controllers

import (
	"../models"
	"../requests"
	"../utils"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const apiKey = "secretApiKey"
const apiKeyNotFound = "Not found"

type TinderService struct {
	Db *sql.DB
}

// CreateUser randomly creates users and adds them to the DB
// Returns a CreateUserResponse
func (s *TinderService) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := utils.GenerateUser()
	token := utils.Base64Encode(user.Email, user.Password)
	user.Token = token
	statement, _ := s.Db.Prepare("INSERT OR IGNORE INTO users (name, email, password, gender, age, token) VALUES (?, ?, ?, ?, ?, ?)")
	res, err := statement.Exec(user.Name, user.Email, user.Password, user.Gender, user.Age, token)
	utils.Must(err)
	id, err := res.LastInsertId()
	utils.Must(err)
	user.Id = int(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := requests.CreateUserResponse{
		User: user,
	}
	json.NewEncoder(w).Encode(resp)
	return
}

// Profiles authenticates the request using an api key, fetches a a list of
// potential matching profiles and returns a ProfilesResponse to the user
// The profiles being returned exclude own user profile or any profiles which
// have been swiped by the user.
func (s *TinderService) Profiles(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(mux.Vars(r)["id"])
	utils.Must(err)

	// Decode user login
	var payload requests.ApiKeyBody
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if payload.ApiKey != apiKey {
		requests.HandleRequestForbidden(w, r)
		return
	}

	// Fetch user id
	var user models.User
	row := s.Db.QueryRow("SELECT id, name, email, password, gender, age FROM users WHERE id=?", userId)
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Gender, &user.Age); err == sql.ErrNoRows {
		requests.HandleRequestNotFound(w, r)
		return
	}

	// Fetch all other profiles except own user profile and already swiped profiles
	var profile models.User
	var profiles []models.User
	rows, _ := s.Db.Query("SELECT id, name, email, password, gender, age FROM users WHERE id NOT IN (SELECT potential_match_id FROM swipes WHERE user_id IS ?) AND id IS NOT ?", user.Id, user.Id)
	for rows.Next() {
		rows.Scan(&profile.Id, &profile.Name, &profile.Email, &profile.Password, &profile.Gender, &profile.Age)
		profiles = append(profiles, profile)
	}

	resp := requests.ProfilesResponse{
		PotentialMatches: profiles,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
	return

}

// upsertUserSwipe
func (s *TinderService) upsertUserSwipe(swipe models.Swipe) {
	var potentialMatch models.Swipe
	row := s.Db.QueryRow("SELECT * FROM swipes WHERE user_id=? AND potential_match_id=?", swipe.UserId, swipe.PotentialMatchId)
	if err := row.Scan(&potentialMatch.Id, &potentialMatch.UserId, &potentialMatch.PotentialMatchId, &potentialMatch.Preference, &potentialMatch.Match); err == sql.ErrNoRows {
		selectStatement, _ := s.Db.Prepare("INSERT INTO swipes (user_id, potential_match_id, preference, match) VALUES (?, ?, ?, ?)")
		_, err := selectStatement.Exec(swipe.UserId, swipe.PotentialMatchId, swipe.Preference, 0)
		utils.Must(err)
	} else {
		deleteStatement, _ := s.Db.Prepare("DELETE FROM swipes WHERE id=?")
		_, err := deleteStatement.Exec(potentialMatch.Id)
		utils.Must(err)

		insertStatement, _ := s.Db.Prepare("INSERT INTO swipes (user_id, potential_match_id, preference, match) VALUES (?, ?, ?, ?)")
		_, err = insertStatement.Exec(swipe.UserId, swipe.PotentialMatchId, swipe.Preference, 0)
		utils.Must(err)
	}
}

// Swipe authenticates the request using an API key.
// When a user swipes, this data is saved in the DB.
// If the potential match has swiped Yes on this user then it will return a match.
func (s *TinderService) Swipe(w http.ResponseWriter, r *http.Request) {
	// Decode user swipe
	var swipe models.Swipe
	err := json.NewDecoder(r.Body).Decode(&swipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if swipe.ApiKey != apiKey {
		requests.HandleRequestForbidden(w, r)
		return
	}

	// Store user swipe
	s.upsertUserSwipe(swipe)

	var resp requests.SwipeResponse
	// Fetch profile of potential match and check preference
	var potentialMatch models.Swipe
	row := s.Db.QueryRow("SELECT * FROM swipes WHERE user_id=? AND potential_match_id=?", swipe.PotentialMatchId, swipe.UserId)
	if err := row.Scan(&potentialMatch.Id, &potentialMatch.UserId, &potentialMatch.PotentialMatchId, &potentialMatch.Preference, &potentialMatch.Match); err == sql.ErrNoRows {
		resp.Match = false
	} else if swipe.Preference == true && potentialMatch.Preference == true {
		swipe.Match = true
		potentialMatch.Match = true
		resp.Match = true
		s.upsertUserSwipe(swipe)
		s.upsertUserSwipe(potentialMatch)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
	return
}

// Login generates a hash based on the user provided credentials (email and password)
// If the hash is matching the one stored in the DB for the user which is
// registered under this email and password, then Login returns an API key.
// API Key is used for authenticating further requests.
func (s *TinderService) Login(w http.ResponseWriter, r *http.Request) {
	// Decode user login
	var login requests.Login
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var storedToken string

	generatedToken := utils.Base64Encode(login.Email, login.Password)

	row := s.Db.QueryRow("SELECT token FROM users WHERE email=? AND password=?", login.Email, login.Password)
	err = row.Scan(&storedToken)
	if err != nil {
		requests.HandleRequestNotFound(w, r)
		return
	}

	var response requests.LoginResponse
	if storedToken == generatedToken {
		response.ApiKey = apiKey
	} else {
		response.ApiKey = apiKeyNotFound
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
	return
}
