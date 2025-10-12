package main

import (
	"net/http"
	"encoding/json"
	"strings"
)

var profanity = [...]string{"kerfuffle", "sharbert", "fornax"}

func validateChirp(w http.ResponseWriter, req *http.Request){
	type parameters struct {
		Body string `json:"body"`
	}	
	type returnValues struct {
		Valid bool `json:"valid"`
		CleanBody string `json:"cleaned_body"`
	}
	
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	} 

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}
	cleanBody := replaceProfanity(params.Body, badWords)

	respondWithJSON(w, http.StatusOK, returnValues{
		Valid: true,
		CleanBody: cleanBody,
	})
}

func replaceProfanity(chirp string, badWords map[string]struct{}) string {
	words := strings.Fields(chirp)

	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if _, ok := badWords[lowerWord]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")	
}

