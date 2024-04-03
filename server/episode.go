package server

import (
	"os"
	"path/filepath"
	"sort"

	. "github.com/d4mr/adcast/podcast"
)

func ScanMediaDir(mediaDir string) (*Podcast, error) {
	// Open the directory
	dir, err := os.Open(mediaDir)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// Read the directory
	files, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	// For each directory in the media directory
	// 	- read the files in the directory
	// 	- create a new PodcastEpisode
	// 	- set the title to the directory name
	// 	- set the partitions to the files in the partitions directory
	// 	- set the ads to the files in the ads directory
	// 	- append the PodcastEpisode to the episodes slice
	// return the episodes slice

	episodes := make(map[string]PodcastEpisode)

	for _, file := range files {
		if file.IsDir() {
			partitionsDir, err := os.Open(mediaDir + "/" + file.Name() + "/partitions")
			if err != nil {
				return nil, err
			}
			defer partitionsDir.Close()

			partitionFiles, err := partitionsDir.Readdir(0)
			if err != nil {
				return nil, err
			}

			adsDir, err := os.Open(mediaDir + "/" + file.Name() + "/ads")
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

			episodes[file.Name()] = PodcastEpisode{
				Title:      file.Name(),
				Partitions: partitions,
				Ads:        ads,
			}
		}
	}

	absMediaDir, err := filepath.Abs(mediaDir)
	if err != nil {
		return nil, err
	}

	return &Podcast{MediaDir: absMediaDir, PodcastEpisodes: episodes}, nil
}
