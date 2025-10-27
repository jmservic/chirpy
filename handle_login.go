package main

import (
	"net/http"
	"encoding/json"
	"github.com/jmservic/chirpy/internal/auth"
	"github.com/jmservic/chirpy/internal/database"
	"time"
)

func (cfg apiConfig) handlerLogin(w http.ResponseWriter, res *http.Request) {
	params := struct{
		Email string `json:"email"`
		Password string `json:"password"`
	//	ExpiresInSecs int `json:"expires_in_seconds,omitempty"`
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
	expirationTime := time.Hour
/*	if params.ExpiresInSecs > 0 && params.ExpiresInSecs < 3600 {
		expirationTime = time.Duration(params.ExpiresInSecs) * time.Second
	} */
	
	token, err := auth.MakeJWT(user.ID, cfg.secret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating JWT token", err)
		return
	}
	refreshToken, _ := auth.MakeRefreshToken()
	err = cfg.db.StoreRefreshToken(res.Context(), database.StoreRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error storing refresh token", err)
		return
	}

	rtnVals := struct{
		UserResources
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{
		UserResources: UserResources{
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		},
		Token: token,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, http.StatusOK, rtnVals)
}
