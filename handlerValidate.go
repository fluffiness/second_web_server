package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		Body string `json:"body"`
	}

	type Response struct {
		CleanedBody string `json:"cleaned_body"`
	}

	type ResponseError struct {
		Error string `json:"error"`
	}

	var reqBody Body
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Printf("error: unable to decoder request body. %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(reqBody.Body) > 140 {
		var resBody ResponseError
		resBody.Error = "chirp exceeds character limit 140"
		resData, err := json.Marshal(resBody)
		if err != nil {
			log.Print("error: unable to JSON marshal response body")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resData)
		return
	}

	profane := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(reqBody.Body, " ")
	cleanedWords := make([]string, len(words))

	for i, word := range words {
		if slices.Contains(profane, strings.ToLower(word)) {
			cleanedWords[i] = "****"
			continue
		}
		cleanedWords[i] = words[i]
	}

	var resBody Response
	resBody.CleanedBody = strings.Join(cleanedWords, " ")
	resData, err := json.Marshal(resBody)
	if err != nil {
		log.Print("error: unable to JSON marshal response body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resData)
}
