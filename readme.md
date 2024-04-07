adcast is a simple golang cli for embedding targeted ads into podcast media files
```
                $$\                                $$\     
                $$ |                               $$ |    
 $$$$$$\   $$$$$$$ | $$$$$$$\ $$$$$$\   $$$$$$$\ $$$$$$\   
 \____$$\ $$  __$$ |$$  _____|\____$$\ $$  _____|\_$$  _|  
 $$$$$$$ |$$ /  $$ |$$ /      $$$$$$$ |\$$$$$$\    $$ |    
$$  __$$ |$$ |  $$ |$$ |     $$  __$$ | \____$$\   $$ |$$\ 
\$$$$$$$ |\$$$$$$$ |\$$$$$$$\\$$$$$$$ |$$$$$$$  |  \$$$$  |
 \_______| \_______| \_______|\_______|\_______/    \____/ 	
```

> [!NOTE]
> This is experimental software, currently requires preprocessing of media files to work. Only supports video files at the moment. This is a proof of concept, and not intended for production use.

## Installation
```bash
go install github.com/d4mr/adcast
```

## Usage
First create a `media` directory with the following structure
```
media
├── progressive
│   ├── [episode_name]
│   │   ├── ads
│   │   │   ├── ad_1.mp4
│   │   │   ├── ad_2.mp4
│   │   │   └── ...
│   │   └── partitions
│   │       ├── 1.mp4
│   │       ├── 2.mp4
│   │       └── ...
│   └── ...
└── hls
    ├── [some_other_episode]
    │   ├── ads
    │   │   ├── ad1.m3u8 
    │   │   ├── ad2.m3u8 # HLS manifest of the ad
    │   │   └── ...
    │   ├── playlist.m3u8 # HLS manifest of the episode
    │   └── ad_times.txt # File containing the timestamps of the ad breaks
    └── ...
```
To start the server run
```bash
adcast server -m [path to media directory]
```
Then navigate to `localhost:8080/hls/[episode_name]` or `localhost:8080/progressive/[episode_name]` to view the episode with ads embedded.

An example media directory is provided in the `example` directory. Content is from the [Creative Commons](https://creativecommons.org/) licensed [Sprite Fight](https://studio.blender.org/films/sprite-fright/) Blender Foundation movie.

## How it works
### HLS Mode
In HLS mode, a dynamic HLS manifest is generated. The media directory contains a `playlist.m3u8` file, which is the HLS manifest of the episode, and a directory of ads, each containing an HLS manifest of the ad.

Simple FFMPEG example for creating the episode manifest:
```bash
ffmpeg -i original.mp4 -c copy -f hls -hls_time 10 -hls_list_size 0 -hls_segment_filename ./output%03d.ts output.m3u8
```
The server reads the episode manifest, and replaces the ad tags with the ad HLS manifest.
This is a lightweight operation, and allows the `ts` files to be served over a CDN. This is the preferred mode of operation.

#### Limitations
Ad insertion can only happen between segments. This means that the HLS manifest must be generated with a segment duration that is short enough to allow for ad insertion. This is a limitation of the HLS protocol.

### Progressive Mode
Progressive mode works over MP4 and other progressive media files. It is preferred only for compatibility with non HLS clients.
Ad order is randomised and ads are inserted between partitions. Partitions are created using `ffmpeg`, with a command like:
```bash
ffmpeg -i original.mp4 -c copy -map 0 -f segment -segment_times 30.0,150.0,270.0,390.0,510.0 -reset_timestamps 1 ./partitions/output%03d.mp4
```
Partitioning like this is fast, and does not require re-encoding the video.

#### Limitations
Partitioning happens only "roughly" at the timestamps specified. Partitioning is not frame accurate, and may not be accurate to the second. This is because ffmpeg will only split the video at keyframes, which may not align with the timestamps specified. This is a limitation of the ffmpeg segment muxer.

If frame accurate partitioning is required, the video must be re-encoded with the correct keyframe interval.

Since every video is "unique", videos cannot be served over CDN. This is because the server must be able to insert ads into the video, and the server must be able to serve the video with ads embedded. This is not possible with a static file. This is why HLS mode is preferred.

#### Ad Insertion
Ads are transcoded to the same codec and resolution as the podcast media, then the FFMPEG concat demuxer is used to concatenate the partitions and ads together. This is very fast, and does not require re-encoding the video.

The equivalent command looks something like:
```bash
ffmpeg -f concat -safe 0 -i list.txt -c copy output.mp4
```
Where `list.txt` is a file containing the paths to the partitions and ads in the correct order.

## Functionality
Ads are shuffled and embedded between the paritions. Upon first request, the server sets a `seed` cookie to ensure the same ad order is used for the duration of the episode.
The cookie can be refreshed by navigating to `localhost:8080/clear-cookie`.

## Why
The goal here is to run a server capable of delivering dynamic media files. Then point a podcast feed to the server, and have the server deliver the media files with ads embedded. This way, the podcast feed can remain static, and the server can be used to manage ads and ad insertion.
This is a test setup to see how podast clients handle dynamic media files, and how well this setup can scale.


## Future
- **Support audio files and other media**
  
- **Perform pre-processing on startup** 
  > Currently for progressive media delivery, media must already be partitioned. This should ideally be done on startup automatically, parametrised by a config file.

- **Decouple ad insertion**
  > Currently it just shuffles ads, which is sufficient for a POC but is not useful. A configuration strategy is needed for the service to enquire which ad to insert when it encounters an ad break.