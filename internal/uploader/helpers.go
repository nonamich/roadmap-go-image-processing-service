package uploader

import "strings"

func GetUploadsDir() string {
	return "storage/uploads/"
}

func (file UploadedFile) GetRelativePath() string {
	dir := GetUploadsDir()
	ext := strings.TrimPrefix(file.Mime, "image/")

	return dir + file.Uuid + "." + ext
}
