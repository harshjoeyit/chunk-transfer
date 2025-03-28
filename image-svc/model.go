package imagesvc

type ImageFile struct {
	// File size in bytes
	Size        int64
	B64         string
	ContentType string
}
