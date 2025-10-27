package main

import (
	"net/http"
	"github.com/jmservic/chirpy/internal/auth"
	"github.com/jmservic/chirpy/internal/database"
	"encoding/json"
)

func (cfg apiConfig) handlerUpdateUser(w http.ResponseWriter, req *http.Request) {
	authToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error getting JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(authToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error validating JWT", err)
		return
	}

	params := struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}
	
	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing new password", err)
		return
	}
	updatedUserInfo, err := cfg.db.UpdateUser(req.Context(), database.UpdateUserParams{
		Email: params.Email,
		HashedPassword: hashed_password,
		ID: userID,})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating users", err)
		return		
	}

	respondWithJSON(w, http.StatusOK, UserResources{
		Id: updatedUserInfo.ID,
		Email:  updatedUserInfo.Email,
		CreatedAt: updatedUserInfo.CreatedAt,
		UpdatedAt: updatedUserInfo.UpdatedAt,
	})
}
