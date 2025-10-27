package main

import (
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"time"
	"github.com/jmservic/chirpy/internal/auth"
	"github.com/jmservic/chirpy/internal/database"
)

type UserResources struct{
	Id uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	ChirpyRed bool `json:"is_chirpy_red"`
}

func (cfg apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
	params := struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}{}	
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	if params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "password cannot be empty.", nil)
		return
	}
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing the password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hash,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating a new user", err)
		return
	}

	rtnVals := UserResources{
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		ChirpyRed: user.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusCreated, rtnVals)
}
