package handlers

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
)

func Image_handler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving the file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)

	if err != nil {
		http.Error(w, "Error decoding image: "+err.Error(), http.StatusBadRequest)
		return
	}
	result := image.NewRGBA(img.Bounds())
	result.At(0, 0)

	fmt.Fprintf(w, "File %s was received successfully!", handler.Filename)
}
