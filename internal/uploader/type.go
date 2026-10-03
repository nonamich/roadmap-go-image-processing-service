package uploader

type UploadedFile struct {
	ID       int64        `json:"id"`
	Uuid     string       `json:"uuid"`
	Mime     string       `json:"mime"`
	Metadata FileMetadata `json:"metadata"`
	UserId   int64        `json:"user_id"`
}

type FileMetadata struct {
	Size     uint64 `json:"size"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Original string `json:"original"`
}
