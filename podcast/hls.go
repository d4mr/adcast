package podcast

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/grafov/m3u8"
)

type HlsEpisode struct {
	Title   string
	Episode *m3u8.MediaPlaylist
	Ads     []*m3u8.MediaPlaylist
	AdTimes []int
}

func (h HlsEpisode) String() string {
	var str string
	str += "Title: " + h.Title + "\t"
	str += fmt.Sprintf("%d Ads\t", len(h.Ads))
	str += fmt.Sprintf("Ad Times: %v", h.AdTimes)
	return str
}

func SetHlsEpisodeHeaders(w http.ResponseWriter, episode *HlsEpisode) {
	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s.m3u8\"", episode.Title))
}

func GetHlsEpisode(podcast *Podcast, title string) (*HlsEpisode, error) {
	episode, ok := podcast.HlsEpisodes[title]
	if !ok {
		return nil, fmt.Errorf("episode not found")
	}
	return &episode, nil
}

func GeneratePlaylist(episode *HlsEpisode, randomisedAds []*m3u8.MediaPlaylist) (*m3u8.MediaPlaylist, error) {
	segmentsSize := episode.Episode.Count()
	// fmt.Printf("episode has %d segments\n", segmentsSize)

	for _, ad := range episode.Ads {
		// fmt.Printf("ad %d has %d segments\n", i, ad.Count())
		segmentsSize += ad.Count()
	}
	generatedPlaylist, err := m3u8.NewMediaPlaylist(segmentsSize, segmentsSize)

	generatedPlaylist.MediaType = m3u8.VOD

	if err != nil {
		return nil, fmt.Errorf("error creating playlist: %w", err)
	}

	// go through each segment, and track total duration
	// as soon as total duration exceeds one of the ad times, insert an ad
	// ad insertion consists of adding discontinuity tag, and then adding the ad segment

	// total duration of the episode
	seek := 0.000

	// index of the ad time we are currently looking at
	adIdx := 0

	episodeSegments := episode.Episode.GetAllSegments()

	wasLastSegmentAd := false

	for _, segment := range episodeSegments {
		// if we have an ad to insert
		// fmt.Printf("Current size %d", generatedPlaylist.Count())

		if adIdx < len(episode.AdTimes) && seek >= float64(episode.AdTimes[adIdx]) {
			// fmt.Println("Trying to add ad segments")

			adSegments := randomisedAds[adIdx].GetAllSegments()
			for i, adSegment := range adSegments {
				if i == 0 {
					// insert discontinuity tag before the ad
					adSegment.Discontinuity = true
				}

				err := generatedPlaylist.AppendSegment(adSegment)
				if err != nil {
					return nil, fmt.Errorf("error adding ad segment: %w", err)
				}

				// fmt.Println("Added ad segment")
			}
			// increment the ad index
			adIdx++
			wasLastSegmentAd = true
		}

		if wasLastSegmentAd {
			// insert discontinuity tag after the ad
			segment.Discontinuity = true
			wasLastSegmentAd = false
		}

		// add the segment
		err := generatedPlaylist.AppendSegment(segment)

		if err != nil {
			return nil, fmt.Errorf("error adding segment: %w", err)
		}

		// fmt.Println("Added a segment")
		// increment the total duration
		seek += segment.Duration
	}

	generatedPlaylist.Close()

	return generatedPlaylist, nil
}

func ServeManifest(w http.ResponseWriter, manifest *bytes.Buffer) {
	w.Write(manifest.Bytes())
}
