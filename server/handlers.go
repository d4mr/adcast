package server

import (
	"net/http"

	. "github.com/d4mr/adcast/podcast"
)

func ServeEpisode(podcast *Podcast) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.PathValue("title")
		episode, err := GetEpisode(podcast, title)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		SetEpisodeHeaders(w, episode)

		randomisedAds, randSeed := GetRandomizedAds(episode, r)
		SetSeedCookie(w, randSeed)

		videoFiles := GetVideoFiles(podcast, episode, randomisedAds)

		fileListPath, err := CreateFileList(videoFiles)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		videoPath, err := RenderVideo(fileListPath, randSeed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ServeVideo(w, r, videoPath)
	}
}

func ClearSeed(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "seed",
		Value:  "deleted",
		MaxAge: -1,
		Path:   "/",
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Seed cleared"))
}
