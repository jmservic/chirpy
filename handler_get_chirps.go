package main

import (
	"net/http"
	"github.com/google/uuid"
	"sort"
	"time"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request){
	dbChirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Probably getting chirps from the database", err)
		return
	}

	authorID := uuid.Nil
	authorIDStr := req.FormValue("author_id") //Can also use req.URL.Query().Get("author_id")
	if authorIDStr != "" {
		authorID, err = uuid.Parse(authorIDStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}
	}
	chirps := make([]Chirp, 0, len(dbChirps))
	for _, dbChirp := range dbChirps {
		if authorID != uuid.Nil && dbChirp.UserID != authorID {
			continue
		}
		temp := Chirp{
			ID: dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body: dbChirp.Body,
			UserID: dbChirp.UserID,
		}
		chirps = append(chirps, temp)
	}

	sortVal := req.URL.Query().Get("sort")
	if sortVal == "" {
		sortVal = "asc"
	}
	if sortVal == "asc" {
		sort.Slice(chirps, func(i, j int) bool {return chirps[i].CreatedAt.Sub(chirps[j].CreatedAt) < time.Duration(0)})
	} else if sortVal == "desc" {
		sort.Slice(chirps, func(i, j int) bool {return chirps[i].CreatedAt.Sub(chirps[j].CreatedAt) > time.Duration(0)})
	}


	respondWithJSON(w, http.StatusOK, chirps)
}
