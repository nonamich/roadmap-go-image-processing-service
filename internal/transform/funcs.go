package transform

import (
	"image"
	"nonamich/image-processing-service/internal/uploader"
	"os"

	"github.com/disintegration/imaging"
)

func Transform(transforms Transforms, uploadedFileData uploader.UploadedFile) (uploader.UploadedFile, error) {
	path := uploadedFileData.GetRelativePath()
	src, err := imaging.Open(path)

	if err != nil {
		return uploadedFileData, err
	}

	if transforms.Resize != nil {
		src = Resize(src, *transforms.Resize)
	}

	err = imaging.Save(src, path)

	if err != nil {
		return uploadedFileData, err
	}

	fileInfo, _ := os.Stat(path)

	bounds := src.Bounds()
	newWidth := bounds.Dx()
	newHeight := bounds.Dy()
	sizeBytes := fileInfo.Size()

	uploadedFileData.Metadata.Width = newWidth
	uploadedFileData.Metadata.Height = newHeight
	uploadedFileData.Metadata.Size = uint64(sizeBytes)

	return uploadedFileData, nil
}

func Resize(src image.Image, resize ResizeTransform) image.Image {
	return imaging.Resize(src, resize.Width, resize.Height, imaging.Lanczos)
}
