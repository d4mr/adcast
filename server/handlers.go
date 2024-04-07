package server

import (
	"fmt"
	"net/http"

	"github.com/d4mr/adcast/podcast"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Crendentials", "true")
}

func ServeProgressiveEpisode(p *podcast.Podcast) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.PathValue("title")
		episode, err := podcast.GetProgressiveEpisode(p, title)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		podcast.SetProgressiveEpisodeHeaders(w, episode)

		randomisedAds, randSeed := podcast.GetRandomizedAds(episode.Ads, r)
		podcast.SetSeedCookie(w, randSeed)

		videoFiles := podcast.GetVideoFiles(p, episode, randomisedAds)

		fileListPath, err := podcast.CreateFileList(videoFiles)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		videoPath, err := podcast.RenderVideo(fileListPath, randSeed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		enableCors(&w)
		podcast.ServeVideo(w, r, videoPath)
	}
}

func ServeHlsEpisode(p *podcast.Podcast) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.PathValue("title")
		episode, err := podcast.GetHlsEpisode(p, title)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			fmt.Println(err)
			return
		}

		podcast.SetHlsEpisodeHeaders(w, episode)

		randomisedAds, randSeed := podcast.GetRandomizedAds(episode.Ads, r)
		podcast.SetSeedCookie(w, randSeed)

		playlist, err := podcast.GeneratePlaylist(episode, randomisedAds)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		enableCors(&w)
		podcast.ServeManifest(w, playlist.Encode())
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
