package imagesvc

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	b64 "encoding/base64"
)

// GetThumbnailDataChunks returns a channel which on read gives the
// thumbnail chunks
func GetThumbnailDataChunks(pathList []string) <-chan string {
	N := len(pathList)

	// Buffered channels with size equal to # of thumnail requested
	ch := make(chan string, N)

	// Since we want return the channel without blocking,
	// start a go-routine to fetch the thumbnail chunks and write
	// those to the channel ch
	go func() {
		var wg sync.WaitGroup
		wg.Add(N)

		for i, path := range pathList {
			// Start a go-routine to get thumbnail chunk denoted by (i, path)
			go func() {
				defer wg.Done()

				chunk, err := GetThumbnailChunk(i, path)
				if err != nil {
					log.Printf("error getting thumbnail chunk for i: %d, path: %s, err: %v", i, path, err)
					return
				}

				// send chunk to the channel
				ch <- chunk
			}()
		}

		// Wait till all go routines are completed
		wg.Wait()

		// All thumbnail chunks sent over the channel.
		// Now, we can safely close the channel
		close(ch)
	}()

	return ch
}

// getThumbnailChunk constructs a chunk for image denoted by path
// which can be written as response to http request
func GetThumbnailChunk(index int, path string) (string, error) {
	// Get the image file content as base64
	f, err := GetThumbanail(path)
	if err != nil {
		return "", err
		// error getting image b64 data for path: %s, error: %v", path,
	}

	// Chunk format: <index>:data:<content-type>;<base64>
	// Example: 0:data:image/png;iVBORw0KGgoAAAANSUhE....
	chunkPrefix := fmt.Sprintf("%d:data:%s;", index, f.ContentType)

	// Write chunk data
	chunk := fmt.Sprintf("%s%s", chunkPrefix, f.B64)

	return chunk, nil
}

// getThumbanail reads the file contents at path and returns as
// base64 with size and content-type
func GetThumbanail(path string) (*ImageFile, error) {
	fpath := filepath.Join(strings.Split(path, "/")...)

	fInfo, err := os.Stat(fpath)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	size := fInfo.Size()
	log.Printf("File size: %d bytes", size)

	data, err := os.ReadFile(fpath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	encoded := b64.StdEncoding.EncodeToString([]byte(data))
	// log.Printf("File b64 encoded str %s", encoded[0:20])

	contentType := http.DetectContentType(data[:512])
	log.Printf("Content type: %s", contentType)

	f := &ImageFile{
		Size:        size,
		B64:         encoded,
		ContentType: contentType,
	}

	return f, nil
}
