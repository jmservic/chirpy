package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, res *http.Request){
	dbChirps, err := cfg.db.GetChirps(res.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Probably getting chirps from the database", err)
		return
	}
	chirps := make([]Chirp, 0, len(dbChirps))
	for _, dbChirp := range dbChirps {
		temp := Chirp{
			ID: dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body: dbChirp.Body,
			UserID: dbChirp.UserID,
		}
		chirps = append(chirps, temp)
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
