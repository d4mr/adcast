package podcast

import (
	"fmt"
	"net/http"
	"path/filepath"
)

type Podcast struct {
	MediaDir        string
	PodcastEpisodes PodcastEpisodes
}

type PodcastEpisodes map[string]PodcastEpisode

type PodcastEpisode struct {
	Title      string
	Partitions []string
	Ads        []string
}

func GetEpisode(podcast *Podcast, title string) (*PodcastEpisode, error) {
	episode, ok := podcast.PodcastEpisodes[title]
	if !ok {
		return nil, fmt.Errorf("episode not found")
	}
	return &episode, nil
}

func SetEpisodeHeaders(w http.ResponseWriter, episode *PodcastEpisode) {
	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s.mp4\"", episode.Title))
}

func GetVideoFiles(podcast *Podcast, episode *PodcastEpisode, randomisedAds []string) []string {
	var videoFiles []string
	for idx, partition := range episode.Partitions {
		videoFiles = append(videoFiles, filepath.Join(podcast.MediaDir, episode.Title, "partitions", partition))
		if idx < len(episode.Ads) {
			videoFiles = append(videoFiles, filepath.Join(podcast.MediaDir, episode.Title, "ads", randomisedAds[idx]))
		}
	}
	return videoFiles
}
