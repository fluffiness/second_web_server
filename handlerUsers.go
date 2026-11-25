package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type UserReqBody struct {
		Email string `json:"email"`
	}

	var reqBody UserReqBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Printf("error: unable to decoder request body. %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), reqBody.Email)
	if err != nil {
		log.Printf("error: failed to create user in database. %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	newUserData, err := json.Marshal(newUser)
	if err != nil {
		log.Printf("error: unable to JSON marshal response body. %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(newUserData)
}
