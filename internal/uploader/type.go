package uploader

type UploadedFile struct {
	ID       int64
	Uuid     string
	Mime     string
	Metadata FileMetadata
	UserId   int64
}

type FileMetadata struct {
	Size     uint64
	Width    int
	Height   int
	Original string
}
