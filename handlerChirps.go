package main

import (
	"encoding/json"
	"fluffiness/second_web_server/internal/database"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type ChirpReqBody struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	var reqBody ChirpReqBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Printf("error: unable to decoder request body. %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	params := database.CreateChirpParams{
		Body:   reqBody.Body,
		UserID: reqBody.UserID,
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), params)
	if err != nil {
		log.Printf("error: failed to create chirp in database. %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	newChirpData, err := json.Marshal(&newChirp)
	if err != nil {
		log.Printf("error: unable to JSON marshal response body. %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(newChirpData)

}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Printf("error: failed to fetch chirps from database. %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonChirps := make([]Chirp, len(chirps))
	for i, c := range chirps {
		jsonChirps[i] = Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		}
	}

	chirpsData, err := json.Marshal(&jsonChirps)
	if err != nil {
		log.Printf("error: unable to JSON marshal response body. %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(chirpsData)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("error: failed to parse provided chirp id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), uuid.UUID(id))
	if err != nil {
		log.Printf("error: failed to fetch chirp from database. %s", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonChrip := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	chirpData, err := json.Marshal(&jsonChrip)
	if err != nil {
		log.Printf("error: unable to JSON marshal response body. %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(chirpData)
}
