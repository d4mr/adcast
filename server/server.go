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

	fmt.Println("Server starting on port 8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
