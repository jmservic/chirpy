package main

import(
	"net/http"
	"github.com/google/uuid"
	"github.com/jmservic/chirpy/internal/auth"
)

func (cfg apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {
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
	
	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error parsing chirp ID path value", err)
		return
	}
	chirp, err := cfg.db.GetChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error getting chirp", err)
	}
	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "UserID of JWT and chirp do not match", nil)
		return
	}

	err = cfg.db.DeleteChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting chirp", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
