package main

import (
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

func (cfg apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
	params := struct{
		Email string `json:"email"`
	}{}	
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating a new user", err)
		return
	}

	rtnVals := struct{
		Id uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
	}{
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}
	respondWithJSON(w, http.StatusCreated, rtnVals)
}
