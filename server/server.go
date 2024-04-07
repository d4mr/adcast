package server

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/d4mr/adcast/podcast"
)

func StartServer(mediaDir string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthCheck)

	p, err := ScanMediaDir(filepath.Join(mediaDir))

	if err != nil {
		fmt.Println("Error scanning media directory: ", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded %d progressive episodes:\n", len(p.ProgressiveEpisodes))

	for _, episode := range p.ProgressiveEpisodes {
		fmt.Println(episode)
	}

	fmt.Printf("\nLoaded %d hls episodes:\n", len(p.HlsEpisodes))

	for _, episode := range p.HlsEpisodes {
		fmt.Println(episode)
	}

	mux.HandleFunc("GET /progressive/episodes/{title}", ServeProgressiveEpisode(p))
	mux.HandleFunc("GET /hls/episodes/{title}", ServeHlsEpisode(p))

	mux.HandleFunc("/clear-seed", ClearSeed)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("\n\nServer starting on port ", port)

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
