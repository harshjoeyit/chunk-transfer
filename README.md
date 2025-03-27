# Chunk Transfer Encoding Prototype

This project demonstrates chunk transfer encoding for sending base64 encoded images using a Gin server.

## Features

*   Sends base64 encoded images in chunks.
*   Demonstrates how to handle chunked responses in JavaScript.
*   Supports fetching multiple thumbnails in a single request.

## Usage

1.  Place image files in the `./images` directory.
2.  Run the Go server: `go run main.go`
3.  Open `index.html` in your browser.

## Endpoints

*   `/data/`: Simple API to demonstrate chunk transfer encoding
*   `/thumbnail-batch?paths=<image_paths>`: Returns base64 encoded images in chunks. `<image_paths>` is a comma-separated list of image paths (e.g., `/images/timg1.png,/images/timg2.png`). Example - `http://localhost:8080/thumbnail-batch?paths=/images/timg1.png,/images/timg2.png,/images/timg3.png`
