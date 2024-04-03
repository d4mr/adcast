package server

import (
	"fmt"
	"net/http"
	"os"
)

func StartServer(mediaDir string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthCheck)

	episodes, err := ScanMediaDir(mediaDir)
	if err != nil {
		fmt.Println("Error scanning media directory: ", err)
		os.Exit(1)
	}

	fmt.Println("Loaded episodes:", episodes)

	mux.HandleFunc("GET /episodes/{title}", ServeEpisode(episodes))
	mux.HandleFunc("/clear-seed", ClearSeed)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server starting on port ", port)

	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
