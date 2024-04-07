package podcast

import (
	"fmt"
	"net/http"
	"path/filepath"
)

type ProgressiveEpisode struct {
	Title      string
	Partitions []string
	Ads        []string
}

func (p ProgressiveEpisode) String() string {
	var str string
	str += fmt.Sprintf("Title: %s\t", p.Title)
	str += fmt.Sprintf("%d Episode Partitions\t", len(p.Partitions))
	str += fmt.Sprintf("%d Ads", len(p.Ads))
	return str
}

func GetProgressiveEpisode(podcast *Podcast, title string) (*ProgressiveEpisode, error) {
	episode, ok := podcast.ProgressiveEpisodes[title]
	if !ok {
		return nil, fmt.Errorf("episode not found")
	}
	return &episode, nil
}

func SetProgressiveEpisodeHeaders(w http.ResponseWriter, episode *ProgressiveEpisode) {
	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s.mp4\"", episode.Title))
}

func GetVideoFiles(podcast *Podcast, episode *ProgressiveEpisode, randomisedAds []string) []string {
	var videoFiles []string
	for idx, partition := range episode.Partitions {
		videoFiles = append(videoFiles, filepath.Join(podcast.MediaDir, episode.Title, "partitions", partition))
		if idx < len(episode.Ads) {
			videoFiles = append(videoFiles, filepath.Join(podcast.MediaDir, episode.Title, "ads", randomisedAds[idx]))
		}
	}
	return videoFiles
}
