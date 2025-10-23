package main

import (
	"net/http"
	"encoding/json"
	"github.com/jmservic/chirpy/internal/auth"
)

func (cfg apiConfig) handlerLogin(w http.ResponseWriter, res *http.Request) {
	params := struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}{}	

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}	

	user, err := cfg.db.GetUserByEmail(res.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error checking hash", err)
		return
	}
	if !match {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return

	}
	rtnVals := UserResources{
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}
	respondWithJSON(w, http.StatusOK, rtnVals)
}
