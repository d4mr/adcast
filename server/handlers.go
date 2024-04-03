package server

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func ServeEpisode(podcast *Podcast) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.PathValue("title")
		episode, ok := podcast.PodcastEpisodes[title]

		if !ok {
			http.Error(w, "Episode not found", http.StatusNotFound)
			return
		}

		// Set the appropriate headers for serving the video
		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s.mp4\"", episode.Title))

		var videoFiles []string

		randomisedAds := make([]string, len(episode.Ads))
		// randonmize the ads
		copy(randomisedAds, episode.Ads)

		randSeed := rand.Int63()
		// read the seed cookie
		seedCookie, err := r.Cookie("seed")
		if err == nil {
			fmt.Println("Seed cookie found:", seedCookie.Value)
			seed, err := strconv.ParseInt(seedCookie.Value, 10, 64)
			if err != nil {
				fmt.Println("Error parsing seed cookie: ", err)
			}
			randSeed = int64(seed)
		} else {
			fmt.Println("No seed cookie found")
		}
		// random src
		rsrc := rand.New(rand.NewSource(randSeed))
		rsrc.Shuffle(len(randomisedAds), func(i, j int) {
			randomisedAds[i], randomisedAds[j] = randomisedAds[j], randomisedAds[i]
		})
		// set seed cookie to be able to reproduce the same random order
		http.SetCookie(w, &http.Cookie{
			Name:  "seed",
			Value: fmt.Sprintf("%d", randSeed),
			Path:  "/",
		})

		for idx, partition := range episode.Partitions {
			videoFiles = append(videoFiles, filepath.Join(podcast.MediaDir, title, "partitions", partition))
			if idx < len(episode.Ads) {
				videoFiles = append(videoFiles, filepath.Join(podcast.MediaDir, title, "ads", randomisedAds[idx]))
			}
		}

		// Create a temporary file to store the list of video files
		tempDir, err := os.MkdirTemp("", "ffmpeg-concat")
		if err != nil {
			http.Error(w, "Error creating temp dir", http.StatusInternalServerError)
			fmt.Println("Error creating temp dir: ", err)
			return
		}
		defer os.RemoveAll(tempDir)

		// Create the file list
		fileListPath := filepath.Join(tempDir, "files.txt")

		var fileListBuffer bytes.Buffer
		for _, file := range videoFiles {
			fileListBuffer.WriteString(fmt.Sprintf("file '%s'\n", file))
		}
		err = os.WriteFile(fileListPath, fileListBuffer.Bytes(), 0644)

		if err != nil {
			http.Error(w, "Error writing temp file", http.StatusInternalServerError)
			fmt.Println("Error writing temp file: ", err)
			return
		}

		tempFile, err := os.CreateTemp("", "concatenated_episode_*.mp4")
		if err != nil {
			fmt.Println("Failed to create temporary file")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer os.Remove(tempFile.Name())

		var errorBuffer bytes.Buffer
		err = ffmpeg.Input(fileListPath, ffmpeg.KwArgs{"f": "concat", "protocol_whitelist": "file,http,https,tcp,tls,crypto,data,pipe", "safe": 0}).
			Output(tempFile.Name(), ffmpeg.KwArgs{"c": "copy", "f": "mp4"}).
			OverWriteOutput().
			WithErrorOutput(&errorBuffer).
			Silent(true).
			Run()

		if err != nil {
			fmt.Println("Error rendering video: ", err)
			fmt.Println("Error output: ", errorBuffer.String())
			http.Error(w, "Error rendering video", http.StatusInternalServerError)
			return
		}

		http.ServeFile(w, r, tempFile.Name())
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
