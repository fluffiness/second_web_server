package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type UserEmail struct {
		Email string `json:"email"`
	}

	var userEmail UserEmail
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userEmail)
	if err != nil {
		log.Printf("error: unable to decoder request body. %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), userEmail.Email)
	newUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	newUserData, err := json.Marshal(newUser)
	if err != nil {
		log.Print("error: unable to JSON marshal response body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(newUserData)
}
