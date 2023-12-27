package handlers

import (
	"image"
	"image/png"
	"math"
	"net/http"
	L "seam-carving/utils"
	"strconv"
)

func Image_handler(w http.ResponseWriter, r *http.Request) {
	img, err := L.GetImageFromRequest(r)
	if err != nil {
		http.Error(w, "Error decoding image: "+err.Error(), http.StatusBadRequest)
		return
	}
	width, err := strconv.Atoi(r.FormValue("width"))
	if err != nil {
		http.Error(w, "Invalid width value", http.StatusBadRequest)
	}
	height, err := strconv.Atoi(r.FormValue("height"))
	if err != nil {
		http.Error(w, "Invalid height value", http.StatusBadRequest)
	}

	energies := make([][]float64, img.Bounds().Dx())
	for i := range energies {
		energies[i] = make([]float64, img.Bounds().Dy())
	}

	result := image.NewRGBA(img.Bounds())
	for i := 0; i < len(energies); i++ {
		for j := 0; j < len(energies[0]); j++ {
			result.Set(i, j, img.At(i, j))
			energies[i][j] = L.ComputeEnergy(img, i, j)
		}
	}

	maxSteps, err := strconv.Atoi(r.FormValue("max_steps"))
	if err != nil {
		maxSteps = 1
	}

	horizontalDynamic := L.GetHorizontalDynamicPrepResult(energies, maxSteps)

	for i := 0; i < img.Bounds().Dy()-height; i++ {
		seam := L.GetHorizontalSeam(horizontalDynamic, maxSteps)
		for j := 0; j < len(seam); j++ {
			horizontalDynamic[j][seam[j]] = math.Inf(1)
		}
		for k := 0; k < len(horizontalDynamic)-1; k++ {
			left, right := 0, 0
			for right < len(horizontalDynamic[k]) {
				if horizontalDynamic[k][right] != math.Inf(1) {
					tempDyn := horizontalDynamic[k][right]
					horizontalDynamic[k][right] = horizontalDynamic[k][left]
					horizontalDynamic[k][left] = tempDyn

					tempResult := result.At(k, right)
					result.Set(k, right, result.At(k, left))
					result.Set(k, left, tempResult)

					left++
				}
				right++
			}
		}
	}

	verticalDynamic := L.GetVerticalDynamicPrepResult(energies, maxSteps)

	for i := 0; i < img.Bounds().Dx()-width; i++ {
		seam := L.GetVerticalSeam(verticalDynamic, maxSteps) // Using vertical seam calculation
		for j := 0; j < len(seam); j++ {
			verticalDynamic[seam[j]][j] = math.Inf(1)
		}
		for k := 0; k < len(verticalDynamic[0])-1; k++ {
			top, bottom := 0, 0
			for bottom < len(verticalDynamic) {
				if verticalDynamic[bottom][k] != math.Inf(1) {
					tempDyn := verticalDynamic[bottom][k]
					verticalDynamic[bottom][k] = verticalDynamic[top][k]
					verticalDynamic[top][k] = tempDyn

					tempResult := result.At(bottom, k)
					result.Set(bottom, k, result.At(top, k))
					result.Set(top, k, tempResult)

					top++
				}
				bottom++
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	png.Encode(w, result.SubImage(image.Rectangle{Min: image.Pt(0, 0), Max: image.Pt(width, height)}))
}
