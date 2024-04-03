package server

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/d4mr/adcast/podcast"
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

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		podcast.CleanupCache()
		os.Exit(0)
	}()

	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
