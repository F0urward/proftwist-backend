package image

import (
	"errors"
	"net/http"
)

var allowedImagesContentTypes = map[string]string{
	"image/png":  "png",
	"image/jpg":  "jpg",
	"image/jpeg": "jpeg",
}

func CheckImageFileContentType(fileContent []byte) (string, error) {
	contentType := http.DetectContentType(fileContent)

	extension, ok := allowedImagesContentTypes[contentType]
	if !ok {
		return "", errors.New("this content type is not allowed")
	}

	return extension, nil
}
