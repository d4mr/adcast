package podcast

type Podcast struct {
	MediaDir            string
	ProgressiveEpisodes PodcastEpisodes[ProgressiveEpisode]
	HlsEpisodes         PodcastEpisodes[HlsEpisode]
}

type PodcastEpisodes[T ProgressiveEpisode | HlsEpisode] map[string]T
