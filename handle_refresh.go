package main

import(
	"net/http"
	"github.com/jmservic/chirpy/internal/auth"
	"database/sql"
	"time"
)

func (cfg apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)	
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error getting refresh token", err)
		return
	}
	//When checking the expiration, check for the error sql.ErrNoRows, if we get that return 401 else it is a server error.
	isExpired, err := cfg.db.RefreshTokenExpired(req.Context(), refreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
			return
		} else {

			respondWithError(w, http.StatusInternalServerError, "Error checking refresh token", err)
			return
		}	
	}
	if isExpired.Bool {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", nil)
		return
	}

	userId, err := cfg.db.GetUserFromRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting user id", err)
		return
	}

	token, err := auth.MakeJWT(userId, cfg.secret, time.Hour)

	rtnParams := struct{
		Token string `json:"token"`
	}{
		Token: token,
	}

	respondWithJSON(w, http.StatusOK, rtnParams)
}
