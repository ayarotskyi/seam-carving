package utils

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
)

func GetImageFromRequest(r *http.Request) (image.Image, error) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		return nil, err
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}
