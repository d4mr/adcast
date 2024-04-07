package server

import (
	"fmt"
	"net/http"

	. "github.com/d4mr/adcast/podcast"
)

func ServeProgressiveEpisode(podcast *Podcast) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.PathValue("title")
		episode, err := GetProgressiveEpisode(podcast, title)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		SetProgressiveEpisodeHeaders(w, episode)

		randomisedAds, randSeed := GetRandomizedAds(episode.Ads, r)
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

func ServeHlsEpisode(podcast *Podcast) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.PathValue("title")
		episode, err := GetHlsEpisode(podcast, title)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			fmt.Println(err)
			return
		}

		SetHlsEpisodeHeaders(w, episode)

		randomisedAds, randSeed := GetRandomizedAds(episode.Ads, r)
		SetSeedCookie(w, randSeed)

		playlist, err := GeneratePlaylist(episode, randomisedAds)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		ServeManifest(w, playlist.Encode())
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
