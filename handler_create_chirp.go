package main

import (
	"net/http"
	"encoding/json"
	"github.com/jmservic/chirpy/internal/database"
	"github.com/google/uuid"
	"time"
	"github.com/jmservic/chirpy/internal/auth"
)

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error getting JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error validating JWT", err)
		return
	}
	params := struct {
		Body string `json:"body"`
	}{}

	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	} 

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{Body: params.Body, UserID: userID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating the chirp", err)
		return
	}
	rtnVals := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, rtnVals)
}
