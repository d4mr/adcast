package server

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/d4mr/adcast/podcast"
	"github.com/grafov/m3u8"
)

func ScanMediaDir(mediaDir string) (*podcast.Podcast, error) {

	progressiveEpisodes, err := ScanProgressiveMediaDir(filepath.Join(mediaDir, "progressive"))
	if err != nil {
		return nil, fmt.Errorf("error scanning progressive media dir: %w", err)
	}

	hlsEpisodes, err := ScanHlsMediaDir(filepath.Join(mediaDir, "hls"))
	if err != nil {
		return nil, fmt.Errorf("error scanning hls media dir: %w", err)
	}

	absMediaDir, err := filepath.Abs(mediaDir)

	if err != nil {
		return nil, err
	}

	return &podcast.Podcast{
		MediaDir:            absMediaDir,
		ProgressiveEpisodes: *progressiveEpisodes,
		HlsEpisodes:         *hlsEpisodes,
	}, nil
}

func ScanProgressiveMediaDir(progressiveMediaDir string) (*podcast.PodcastEpisodes[podcast.ProgressiveEpisode], error) {
	// Open the directory
	dir, err := os.Open(progressiveMediaDir)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// Read the directory
	files, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	episodes := make(podcast.PodcastEpisodes[podcast.ProgressiveEpisode])

	for _, file := range files {
		if !file.IsDir() {
			return nil, fmt.Errorf("unexpected file in progressive media dir: %s", file.Name())
		}

		partitionsDir, err := os.Open(filepath.Join(progressiveMediaDir, file.Name(), "partitions"))
		if err != nil {
			return nil, err
		}
		defer partitionsDir.Close()

		partitionFiles, err := partitionsDir.Readdir(0)
		if err != nil {
			return nil, err
		}

		adsDir, err := os.Open(filepath.Join(progressiveMediaDir, file.Name(), "ads"))
		if err != nil {
			return nil, err
		}
		defer adsDir.Close()

		adFiles, err := adsDir.Readdir(0)
		if err != nil {
			return nil, err
		}

		var partitions []string
		for _, partitionFile := range partitionFiles {
			partitions = append(partitions, partitionFile.Name())
		}

		//sort partitions by name
		sort.StringSlice(partitions).Sort()

		var ads []string
		for _, adFile := range adFiles {
			ads = append(ads, adFile.Name())
		}

		episodes[file.Name()] = podcast.ProgressiveEpisode{
			Title:      file.Name(),
			Partitions: partitions,
			Ads:        ads,
		}
	}

	return &episodes, nil
}

func ScanHlsMediaDir(hlsMediaDir string) (*podcast.PodcastEpisodes[podcast.HlsEpisode], error) {
	// Open the directory
	dir, err := os.Open(hlsMediaDir)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// Read the directory
	files, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	episodes := make(podcast.PodcastEpisodes[podcast.HlsEpisode])

	for _, file := range files {
		if !file.IsDir() {
			return nil, fmt.Errorf("unexpected file in hls media dir: %s", file.Name())
		}

		episodeManifestFile, err := os.Open(filepath.Join(hlsMediaDir, file.Name(), "playlist.m3u8"))
		if err != nil {
			return nil, err
		}
		defer episodeManifestFile.Close()
		episodeManifestUnknown, episodeManifestListType, err := m3u8.DecodeFrom(bufio.NewReader(episodeManifestFile), true)

		if episodeManifestListType != m3u8.MEDIA {
			return nil, fmt.Errorf("unexpected playlist type for %s: %s", file.Name(), episodeManifestListType)
		}

		episodeManifest := episodeManifestUnknown.(*m3u8.MediaPlaylist)

		if err != nil {
			return nil, fmt.Errorf("error parsing episode manifest for %s: %w", file.Name(), err)
		}

		adManifestsDir, err := os.Open(filepath.Join(hlsMediaDir, file.Name(), "ads"))
		if err != nil {
			return nil, err
		}
		defer adManifestsDir.Close()

		adManifestFiles, err := adManifestsDir.Readdir(0)
		if err != nil {
			return nil, err
		}

		var adManifests []*m3u8.MediaPlaylist
		for _, adManifestFile := range adManifestFiles {
			adManifestFile, err := os.Open(filepath.Join(hlsMediaDir, file.Name(), "ads", adManifestFile.Name()))
			if err != nil {
				return nil, fmt.Errorf("error opening ad manifest for %s: %w", file.Name(), err)
			}

			adManifestUnknown, adManifestListType, err := m3u8.DecodeFrom(bufio.NewReader(adManifestFile), true)

			if adManifestListType != m3u8.MEDIA {
				return nil, fmt.Errorf("unexpected playlist type for %s: %s", file.Name(), adManifestListType)
			}
			adManifest := adManifestUnknown.(*m3u8.MediaPlaylist)
			adManifests = append(adManifests, adManifest)
		}

		// read ad_times.txt
		adTimesFile, err := os.Open(filepath.Join(hlsMediaDir, file.Name(), "ad_times.txt"))
		if err != nil {
			return nil, fmt.Errorf("error opening ad_times.txt for %s: %w", file.Name(), err)
		}

		// ad times are newline separated ints
		// load the entire file, and store it into a slice of ints
		var adTimes []int
		scanner := bufio.NewScanner(adTimesFile)
		for scanner.Scan() {
			adTime, err := strconv.Atoi(scanner.Text())
			if err != nil {
				return nil, fmt.Errorf("error parsing ad time for %s: %w", file.Name(), err)
			}
			adTimes = append(adTimes, adTime)
		}

		episodes[file.Name()] = podcast.HlsEpisode{
			Title:   file.Name(),
			Episode: episodeManifest,
			Ads:     adManifests,
			AdTimes: adTimes,
		}

	}

	return &episodes, nil
}
