package handlers

import (
	"image/color"
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
	originalHeight, originalWidth := img.Bounds().Dy(), img.Bounds().Dx()
	width, err := strconv.Atoi(r.FormValue("width"))
	if err != nil {
		http.Error(w, "Invalid width value", http.StatusBadRequest)
	}
	height, err := strconv.Atoi(r.FormValue("height"))
	if err != nil {
		http.Error(w, "Invalid height value", http.StatusBadRequest)
	}

	energies := make([][]float64, originalWidth)
	for i := range energies {
		energies[i] = make([]float64, originalHeight)
	}

	colorMap := make([][]color.Color, originalWidth)
	for i := 0; i < originalWidth; i++ {
		colorMap[i] = make([]color.Color, originalHeight)
	}
	for i := 0; i < originalWidth; i++ {
		for j := 0; j < originalHeight; j++ {
			colorMap[i][j] = img.At(i, j)
			energies[i][j] = L.ComputeEnergy(img, i, j)
		}
	}

	maxSteps, err := strconv.Atoi(r.FormValue("max_steps"))
	if err != nil {
		maxSteps = 1
	}

	horizontalDynamic := L.GetHorizontalDynamicPrepResult(energies, maxSteps)

	for i := 0; i < originalHeight-height; i++ {
		for col, row := range L.GetHorizontalSeam(horizontalDynamic, maxSteps) {
			removeFloat64AtIndex(horizontalDynamic[col], row)
			removeColorAtIndex(colorMap[col], row)
		}
	}

	verticalDynamic := L.GetVerticalDynamicPrepResult(energies, maxSteps)

	for i := 0; i < originalWidth-width; i++ {
		for row, col := range L.GetVerticalSeam(verticalDynamic, maxSteps) {
			verticalDynamic[col][row] = math.Inf(1)
		}
		for k := 0; k < len(verticalDynamic[0])-1; k++ {
			top, bottom := 0, 0
			for bottom < len(verticalDynamic) {
				if verticalDynamic[bottom][k] != math.Inf(1) {
					tempDyn := verticalDynamic[bottom][k]
					verticalDynamic[bottom][k] = verticalDynamic[top][k]
					verticalDynamic[top][k] = tempDyn

					tempResult := colorMap[bottom][k]
					colorMap[bottom][k] = colorMap[top][k]
					colorMap[top][k] = tempResult

					top++
				}
				bottom++
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	png.Encode(w, L.CreateImageFromColorMap(colorMap, width, height))
}

func removeFloat64AtIndex(s []float64, index int) []float64 {
	return append(s[:index], s[index+1:]...)
}
func removeColorAtIndex(s []color.Color, index int) []color.Color {
	return append(s[:index], s[index+1:]...)
}
