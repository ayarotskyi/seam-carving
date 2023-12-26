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

	horizontalDynamic := getDynamicPrepResult(energies, maxSteps)

	for i := 0; i < img.Bounds().Dy()-height; i++ {
		seam := getHorizontalSeam(horizontalDynamic, maxSteps)
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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	png.Encode(w, result.SubImage(image.Rectangle{Min: image.Pt(0, 0), Max: image.Pt(width, height)}))
}

func getHorizontalSeam(dynamic [][]float64, maxStep int) []int {
	result := make([]int, len(dynamic))
	result[len(dynamic)-1] = func() int {
		min, minIndex := math.Inf(1), 0
		for i := 0; i < len(dynamic[len(dynamic)-1]); i++ {
			temp := dynamic[len(dynamic)-1][i]
			if temp < min {
				min = temp
				minIndex = i
			}
		}
		return minIndex
	}()
	for i := len(dynamic) - 2; i >= 0; i-- {
		minIndex, min := 0, math.Inf(1)
		for k := 0; k < (maxStep*2 + 1); k++ {

			tempIndex := int(math.Max(float64((result[i+1]-maxStep+k)%(len(dynamic[i])-1)), 0))
			temp := dynamic[i][tempIndex]
			if temp < min {
				min = temp
				minIndex = tempIndex
			}
		}
		result[i] = minIndex
	}

	return result
}

func getDynamicPrepResult(energies [][]float64, maxStep int) [][]float64 {
	// creating 2-dimentional array filled with +Inf
	dynamic := make([][]float64, len(energies))
	for i := 0; i < len(dynamic); i++ {
		dynamic[i] = make([]float64, len(energies[0]))
		for j := 0; j < len(dynamic[i]); j++ {
			dynamic[i][j] = energies[i][j]
		}
	}

	// using dynamic programming to get min cumulative energy for each point starting from top
	for i := 1; i < len(dynamic); i++ {
		for j := 0; j < len(dynamic[i]); j++ {
			min := math.Inf(1)
			for k := 0; k < (maxStep*2 + 1); k++ {
				min = math.Min(dynamic[i-1][int(math.Max(float64((j-maxStep+k)%(len(dynamic[i-1])-1)), 0))], min)
			}
			dynamic[i][j] = min + energies[i][j]
		}
	}

	return dynamic
}
