package transform

import (
	"nonamich/image-processing-service/internal/uploader"
)

func Transform(transforms Transforms, file uploader.UploadedFile) {
	if transforms.Resize != nil {

	}
}

func Resize(resize ResizeTransform, u uploader.UploadedFile) {
	path := u.GetRelativePath()

	src, err := imaging.Open(path)
	if err != nil {
		return err
	}
}
