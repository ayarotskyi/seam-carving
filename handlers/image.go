package handlers

import (
	"image/color"
	"image/png"
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
		}
	}
	for i := 0; i < originalWidth; i++ {
		for j := 0; j < originalHeight; j++ {
			energies[i][j] = L.CalculateE1(colorMap, i, j)
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
		getAndRemoveHorizontalSeam(horizontalAcummulativeMatrix, colorMap, maxSteps)
	}

	verticalAcummulativeMatrix := L.GetVerticalAcummulativeMatrix(energies, maxSteps)

	for i := 0; i < originalWidth-width; i++ {
		getAndRemoveVerticalSeam(verticalAcummulativeMatrix, colorMap, maxSteps)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	png.Encode(w, L.CreateImageFromColorMap(colorMap, width, height))
}

func getAndRemoveVerticalSeam(verticalAcummulativeMatrix [][]float64, colorMap [][]color.Color, maxSteps int) {
	width := len(verticalAcummulativeMatrix)
	for row, col := range L.GetVerticalSeam(verticalAcummulativeMatrix, maxSteps) {
		for i := col; i < width-2; i++ {
			verticalAcummulativeMatrix[i][row] = verticalAcummulativeMatrix[i+1][row]
			colorMap[i][row] = colorMap[i+1][row]
		}
	}
	verticalAcummulativeMatrix = verticalAcummulativeMatrix[:width-1]
}

func getAndRemoveHorizontalSeam(horizontalAcummulativeMatrix [][]float64, colorMap [][]color.Color, maxSteps int) {
	for col, row := range L.GetHorizontalSeam(horizontalAcummulativeMatrix, maxSteps) {
		horizontalAcummulativeMatrix[col] = append(horizontalAcummulativeMatrix[col][:row], horizontalAcummulativeMatrix[col][row+1:]...)
		colorMap[col] = append(colorMap[col][:row], colorMap[col][row+1:]...)
	}
}
