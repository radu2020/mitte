package models

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
	Token    string `json:"token"`
}

type Swipe struct {
	Id               int    `json:"id"`
	UserId           int    `json:"user_id"`
	PotentialMatchId int    `json:"potential_match_id"`
	Preference       bool   `json:"preference"`
	Match            bool   `json:"match"`
	ApiKey           string `json:"api_key"`
}
