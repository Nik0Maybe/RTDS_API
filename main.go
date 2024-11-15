package main

import (
	"RTDS_API/handlers"
	"RTDS_API/metrics"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func apiKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.RequestsTotal.WithLabelValues(r.URL.Path, r.Method, "200").Inc()
		key := r.Header.Get("X-API-Key")
		if key != "your-secure-api-key" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			metrics.RequestsTotal.WithLabelValues(r.URL.Path, r.Method, "401").Inc()
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/").Subrouter()
	apiRouter.Use(apiKeyAuth)

	apiRouter.HandleFunc("/stream/start", handlers.StartStreamHandler).Methods("POST")
	apiRouter.HandleFunc("/stream/{stream_id}/send", handlers.SendDataHandler).Methods("POST")
	router.HandleFunc("/stream/{stream_id}/results", handlers.GetResultsHandler).Methods("GET")

	metrics.SetupMetrics(router)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
