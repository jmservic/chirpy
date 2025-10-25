package main

import(
	"net/http"
	"encoding/json"
	"github.com/jmservic/chirpy/internal/auth"
	"database/sql"
)

func (cfg apiConfig) handlerRefresh(w http.ResponseWriter, req http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)	
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error getting refresh token", err)
		return
	}
	//When checking the expiration, check for the error sql.ErrNoRows, if we get that return 401 else it is a server error.
}
