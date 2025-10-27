package main

import(
	"net/http"
	"github.com/jmservic/chirpy/internal/auth"
)

func (cfg apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)	
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error getting refresh token", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
	
