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
	if maxSteps > originalHeight || maxSteps > originalWidth {
		if originalHeight > originalWidth {
			maxSteps = originalHeight
		} else {
			maxSteps = originalWidth
		}
	}

	horizontalAcummulativeMatrix := L.GetHorizontalAcummulativeMatrix(energies, maxSteps)

	for i := 0; i < originalHeight-height; i++ {
		for col, row := range L.GetHorizontalSeam(horizontalAcummulativeMatrix, maxSteps) {
			removeFloat64AtIndex(horizontalAcummulativeMatrix[col], row)
			removeColorAtIndex(colorMap[col], row)
		}
	}

	verticalAcummulativeMatrix := L.GetVerticalAcummulativeMatrix(energies, maxSteps)

	for i := 0; i < originalWidth-width; i++ {
		for row, col := range L.GetVerticalSeam(verticalAcummulativeMatrix, maxSteps) {
			verticalAcummulativeMatrix[col][row] = math.Inf(1)
		}
		for k := 0; k < len(verticalAcummulativeMatrix[0])-1; k++ {
			top, bottom := 0, 0
			for bottom < len(verticalAcummulativeMatrix) {
				if verticalAcummulativeMatrix[bottom][k] != math.Inf(1) {
					tempDyn := verticalAcummulativeMatrix[bottom][k]
					verticalAcummulativeMatrix[bottom][k] = verticalAcummulativeMatrix[top][k]
					verticalAcummulativeMatrix[top][k] = tempDyn

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
