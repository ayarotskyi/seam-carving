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

	energies := make([][]uint32, originalWidth)
	for i := range energies {
		energies[i] = make([]uint32, originalHeight)
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
			energies[i][j] = L.ComputeEnergy(colorMap, i, j)
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

func getAndRemoveVerticalSeam(verticalAcummulativeMatrix [][]uint32, colorMap [][]color.Color, maxSteps int) {
	for row, col := range L.GetVerticalSeam(verticalAcummulativeMatrix, maxSteps) {
		verticalAcummulativeMatrix[col][row] = uint32(math.Inf(1))
	}
	for k := 0; k < len(verticalAcummulativeMatrix[0])-1; k++ {
		top, bottom := 0, 0
		for bottom < len(verticalAcummulativeMatrix) {
			if verticalAcummulativeMatrix[bottom][k] != uint32(math.Inf(1)) {
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

func getAndRemoveHorizontalSeam(horizontalAcummulativeMatrix [][]uint32, colorMap [][]color.Color, maxSteps int) {
	for col, row := range L.GetHorizontalSeam(horizontalAcummulativeMatrix, maxSteps) {
		removeUint32AtIndex(horizontalAcummulativeMatrix[col], row)
		removeColorAtIndex(colorMap[col], row)
	}
}

func removeUint32AtIndex(s []uint32, index int) []uint32 {
	return append(s[:index], s[index+1:]...)
}
func removeColorAtIndex(s []color.Color, index int) []color.Color {
	return append(s[:index], s[index+1:]...)
}
