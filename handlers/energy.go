package handlers

import (
	"image"
	"image/color"
	"image/png"
	"net/http"
	L "seam-carving/utils"
)

func Energy_handler(w http.ResponseWriter, r *http.Request) {
	img, err := L.GetImageFromRequest(r)
	if err != nil {
		http.Error(w, "Error decoding image: "+err.Error(), http.StatusBadRequest)
		return
	}

	energies := make([][]float64, img.Bounds().Dx())
	for i := range energies {
		energies[i] = make([]float64, img.Bounds().Dy())
	}

	max := 0.0
	for i := 0; i < len(energies); i++ {
		for j := 0; j < len(energies[0]); j++ {
			energy := L.ComputeEnergy(img, i, j)

			energies[i][j] = energy
			if energy > max {
				max = energy
			}
		}
	}

	result := image.NewRGBA(img.Bounds())
	for i := 0; i < len(energies); i++ {
		for j := 0; j < len(energies[0]); j++ {
			result.Set(i, j, color.Gray16{0xffff - uint16((energies[i][j]/max)*0xffff)})
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	png.Encode(w, result.SubImage(result.Rect))
}
