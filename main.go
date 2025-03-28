package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	imagesvc "github.com/harshjoeyit/chunk-transfer/image-svc"
)

const chunkSize = 1024 // 1KB chunk size

func main() {
	ge := gin.Default()

	ge.Use(CORSMiddleware())

	// API to send base64 encoded images as chunk encoded
	ge.GET("/thumbnail-batch-concurrent", func(c *gin.Context) {
		// Extract image paths from query parameters
		// ?paths=/images/timg1.png,/images/timg2.png,/images/timg3.png
		paths := c.Query("paths")

		if paths == "" {
			c.String(http.StatusBadRequest, "Missing 'paths' query parameter")
			return
		}

		// Split the paths string into a slice of individual paths
		pathList := strings.Split(paths, ",")

		// Print the extracted paths
		log.Println("Extracted paths:")
		for _, path := range pathList {
			log.Println(path)
		}

		// Set headers for chunked transfer encoding
		// c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Type", "text/plain")
		c.Header("Transfer-Encoding", "chunked")

		var err error

		for dataChunk := range imagesvc.GetThumbnailDataChunks(pathList) {
			// Write the data chunk
			_, err = c.Writer.WriteString(dataChunk)
			if err != nil {
				log.Println("Error writing chunk data:", err)
				return
			}

			// Write CRLF (end of line) to denote end of chunk
			_, err = c.Writer.WriteString("\r\n")
			if err != nil {
				log.Println("Error writing CRLF:", err)
				return
			}

			// Flush to ensure data is sent immediately
			c.Writer.Flush()
		}

		// Write the final chunk (zero-size chunk)
		_, err = c.Writer.WriteString("0\r\n\r\n")
		if err != nil {
			log.Println("Error writing final chunk:", err)
			return
		}

		log.Println("All files sent")
	})

	// API to send base64 encoded images as chunk encoded
	ge.GET("/thumbnail-batch-blocking", func(c *gin.Context) {
		// Extract image paths from query parameters
		// ?paths=/images/timg1.png,/images/timg2.png,/images/timg3.png
		paths := c.Query("paths")

		if paths == "" {
			c.String(http.StatusBadRequest, "Missing 'paths' query parameter")
			return
		}

		// Split the paths string into a slice of individual paths
		pathList := strings.Split(paths, ",")

		// Print the extracted paths
		log.Println("Extracted paths:")
		for _, path := range pathList {
			log.Println(path)
		}

		// Set headers for chunked transfer encoding
		// c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Type", "text/plain")
		c.Header("Transfer-Encoding", "chunked")

		var f *imagesvc.ImageFile
		var err error

		// Todo: Instead of processing images sequentially - for better concurrency
		// we could use goroutines and channels to read files in separate go routines
		// and write data to a channel. We could then for loop on the channel and
		// send the chunks to client
		for i, path := range pathList {
			// Get the image file content as base64
			f, err = imagesvc.GetThumbanail(path)
			if err != nil {
				log.Printf("error getting image b64 data for path: %s, error: %v", path, err)
				continue
			}

			// Chunk format: <index>:data:<content-type>;<base64>
			// Example: 0:data:image/png;iVBORw0KGgoAAAANSUhE....
			chunkPrefix := fmt.Sprintf("%d:data:%s;", i, f.ContentType)

			// Write chunk data
			chunk := fmt.Sprintf("%s%s", chunkPrefix, f.B64)

			_, err = c.Writer.WriteString(chunk)
			if err != nil {
				log.Println("Error writing chunk data:", err)
				return
			}

			// Write CRLF (end of line) to denote end of chunk
			_, err = c.Writer.WriteString("\r\n")
			if err != nil {
				log.Println("Error writing CRLF:", err)
				return
			}

			// Flush to ensure data is sent immediately
			c.Writer.Flush()
		}

		// Write the final chunk (zero-size chunk)
		_, err = c.Writer.WriteString("0\r\n\r\n")
		if err != nil {
			log.Println("Error writing final chunk:", err)
			return
		}

		log.Println("all files sent")
	})

	// A simple API to demonstrate chunk transfer encoding - response is text strings
	ge.GET("/data", func(c *gin.Context) {
		// Simulate a large data source (e.g., a file)
		data := make([]byte, 5*chunkSize) // 5KB of data
		for i := range data {
			data[i] = byte('A' + (i % 26)) // Fill with repeating letters
		}

		// Set headers for chunked transfer encoding
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Transfer-Encoding", "chunked")

		// Write data in chunks
		for i := 0; i < len(data); i += chunkSize {
			end := min(i+chunkSize, len(data))

			chunk := data[i:end]

			// Write chunk size in hex followed by CRLF
			chunkSizeStr := strconv.FormatInt(int64(len(chunk)), 16) + "\r\n"

			log.Println("chunk size: ", len(chunk))
			log.Println("chunk size in hex: ", chunkSizeStr)

			_, err := c.Writer.WriteString(chunkSizeStr)
			if err != nil {
				log.Println("Error writing chunk size:", err)
				return
			}

			// Write chunk data
			_, err = c.Writer.Write(chunk)
			if err != nil {
				log.Println("Error writing chunk data:", err)
				return
			}

			// Write CRLF (end of line) to denote end of chunk
			_, err = c.Writer.WriteString("\r\n")
			if err != nil {
				log.Println("Error writing CRLF:", err)
				return
			}

			c.Writer.Flush() // Flush to ensure data is sent immediately
		}

		// Write the final chunk (zero-size chunk)
		_, err := c.Writer.WriteString("0\r\n\r\n")
		if err != nil {
			log.Println("Error writing final chunk:", err)
			return
		}
	})

	ge.Run(":8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
