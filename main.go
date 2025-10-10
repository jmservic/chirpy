package main

import (
	"net/http"
	"log"
	"sync/atomic"
	"fmt"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = "8080"
	apiCfg := apiConfig {}
	serveMux := http.NewServeMux()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	serveMux.HandleFunc("/healthz", goodHealth)
	serveMux.HandleFunc("/metrics", apiCfg.hitsMetric)
	serveMux.HandleFunc("/reset", apiCfg.resetMetrics)
	server := http.Server{
		Handler: serveMux,
		Addr: ":" + port,
	}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) hitsMetric(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d",cfg.fileserverHits.Load())))
}

