package image

type Image struct {
	ID       int64
	uuid     string
	metadata ImageMetadata
}

type ImageMetadata struct {
	size         uint64
	width        int
	height       int
	originalName string
}
