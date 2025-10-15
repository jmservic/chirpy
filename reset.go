package main

import "net/http"

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	cfg.fileserverHits.Store(0)
	err := cfg.db.DeleteUsers(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating a new user", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0. All Users dropped."))
}
