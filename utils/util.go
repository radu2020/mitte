package utils

import (
	"../models"
	b64 "encoding/base64"
	"math/rand"
)

// Must panics when an error is thrown
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Base64Encode(email, password string) string {
	tokenString := email + ":" + password
	return b64.StdEncoding.EncodeToString([]byte(tokenString))
}

var maleNames = []string{"George", "John", "David", "Fred", "Robert"}
var femaleNames = []string{"Sofia", "Helen", "Natasha", "Giulia", "Laura"}
var gender = []string{"male", "female"}
var ages = []int{25, 26, 28, 30, 32}

func GenerateUser() models.User {
	gender := gender[rand.Intn(len(gender))]
	var name string
	if gender == "male" {
		name = maleNames[rand.Intn(len(maleNames))]
	} else {
		name = femaleNames[rand.Intn(len(maleNames))]
	}
	email := name + "@gmail.com"
	password := name + "s" + "-safe-" + "password"
	age := ages[rand.Intn(len(ages))]

	return models.User{
		Name:     name,
		Email:    email,
		Password: password,
		Gender:   gender,
		Age:      age,
	}
}
