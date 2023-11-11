package handlers

import (
	"fmt"
	"net/http"
	L "seam-carving/utils"
)

func Image_handler(w http.ResponseWriter, r *http.Request) {
	img, err := L.GetImageFromRequest(r)
	if err != nil {
		http.Error(w, "Error decoding image: "+err.Error(), http.StatusBadRequest)
		return
	}

	// width, err := strconv.Atoi(r.FormValue("width"))
	// if err != nil {
	// 	http.Error(w, "Invalid width value", http.StatusBadRequest)
	// }

	// height, err := strconv.Atoi(r.FormValue("height"))
	// if err != nil {
	// 	http.Error(w, "Invalid height value", http.StatusBadRequest)
	// }

	fmt.Fprintf(w, "Image size: %d - width, %d - height", img.Bounds().Dx(), img.Bounds().Dy())
}
