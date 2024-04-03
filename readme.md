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
├── [episode_name]
│   ├── paritions
│   │   ├── 1.mp4
│   │   └── 2.mp4
│   └── ads
│       ├── [ad_name].mp4
│       └── [ad_name].mp4
```
To start the server run
```bash
adcast server -m [path to media directory]
```
Then navigate to `localhost:8080/episodes/[episode_name]` to view the episode with ads embedded.

## How it works
Ad order is randomised and ads are inserted between partitions. Partitions are created using `ffmpeg`, with a command like:
```bash
ffmpeg -i original.mp4 -c copy -map 0 -f segment -segment_times 30.0,150.0,270.0,390.0,510.0 -reset_timestamps 1 ./partitions/output%03d.mp4
```
Partitioning like this is fast, and does not require re-encoding the video.

### Limitations
Partitioning happens only "roughly" at the timestamps specified. Partitioning is not frame accurate, and may not be accurate to the second. This is because ffmpeg will only split the video at keyframes, which may not align with the timestamps specified. This is a limitation of the ffmpeg segment muxer.

If frame accurate partitioning is required, the video must be re-encoded with the correct keyframe interval.

### Ad insertion
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
- Support audio files
- Perform pre-processing on the fly
- Support more config driven ad insertion techniques, like querying an ad server for the next ad to insert
- HLS based ad insertion