package main

import (
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
)

func (cfg apiConfig) handlerChirpyRedUpgrade(w http.ResponseWriter, req *http.Request) {
	params := struct{
		Event string `json:"event"`
		Data struct{
			UserID string `json:"user_id"`
		} `json:"data"`
	}{}
	
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to decode request body", err)
	}
	
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	userId, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to parse user_id", err)
		return
	}

	err = cfg.db.UpgradeUserChirpyRed(req.Context(), userId) 
	if err != nil {
		respondWithError(w, http.StatusNotFound, "unable to upgrade user", err)
	}

	w.WriteHeader(http.StatusNoContent)
}
