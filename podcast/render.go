package podcast

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

var (
	cacheDir      = filepath.Join(os.TempDir(), "video-cache")
	cacheMutex    sync.Mutex
	cacheMaxItems = 10
)

func init() {
	os.MkdirAll(cacheDir, os.ModePerm)
}

func getCachedVideoPath(randSeed int64) string {
	return filepath.Join(cacheDir, fmt.Sprintf("video_%d.mp4", randSeed))
}

func CreateFileList(videoFiles []string) (string, error) {
	tempFilesList, err := os.CreateTemp("", "files.txt")
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %v", err)
	}

	var fileListBuffer bytes.Buffer
	for _, file := range videoFiles {
		fileListBuffer.WriteString(fmt.Sprintf("file '%s'\n", file))
	}

	err = os.WriteFile(tempFilesList.Name(), fileListBuffer.Bytes(), 0644)
	if err != nil {
		return "", fmt.Errorf("error writing temp file: %v", err)
	}

	return tempFilesList.Name(), nil
}

func RenderVideo(fileListPath string, randSeed int64) (string, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	cachedVideoPath := getCachedVideoPath(randSeed)

	// Check if the video is already cached
	if _, err := os.Stat(cachedVideoPath); err == nil {
		fmt.Println("Cache hit ", cachedVideoPath)
		return cachedVideoPath, nil
	}

	var errorBuffer bytes.Buffer
	err := ffmpeg.Input(fileListPath, ffmpeg.KwArgs{"f": "concat", "protocol_whitelist": "file,http,https,tcp,tls,crypto,data,pipe", "safe": 0}).
		Output(cachedVideoPath, ffmpeg.KwArgs{"c": "copy", "f": "mp4"}).
		OverWriteOutput().
		WithErrorOutput(&errorBuffer).
		Silent(true).
		Run()

	if err != nil {
		return "", fmt.Errorf("error rendering video: %v, error output: %s", err, errorBuffer.String())
	}

	// Remove the oldest cached video if the cache exceeds the maximum limit
	cacheItems, _ := os.ReadDir(cacheDir)
	if len(cacheItems) > cacheMaxItems {
		oldestVideo := cacheItems[0].Name()
		os.Remove(filepath.Join(cacheDir, oldestVideo))
		fmt.Println("Removed oldest video from cache: ", oldestVideo)
	}

	return cachedVideoPath, nil
}

func CleanupCache() {
	os.RemoveAll(cacheDir)
	fmt.Println("Cache cleaned up")
}

func ServeVideo(w http.ResponseWriter, r *http.Request, videoPath string) {
	http.ServeFile(w, r, videoPath)
}
