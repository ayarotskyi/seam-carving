package utils

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
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

func GetHorizontalSeam(dynamic [][]float64, maxStep int) []int {
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

			tempIndex := (result[i+1] - maxStep + k) % (len(dynamic[i]) - 1)
			if tempIndex < 0 {
				tempIndex = 0
			}
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

func GetHorizontalDynamicPrepResult(energies [][]float64, maxStep int) [][]float64 {
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
				index := (j - maxStep + k) % (len(dynamic[i-1]) - 1)
				if index < 0 {
					index = 0
				}
				min = math.Min(dynamic[i-1][index], min)
			}
			dynamic[i][j] = min + energies[i][j]
		}
	}

	return dynamic
}

func GetVerticalSeam(dynamic [][]float64, maxStep int) []int {
	result := make([]int, len(dynamic[0]))
	result[len(dynamic[0])-1] = func() int {
		min, minIndex := math.Inf(1), 0
		for i := 0; i < len(dynamic); i++ {
			temp := dynamic[i][len(dynamic[0])-1]
			if temp < min {
				min = temp
				minIndex = i
			}
		}
		return minIndex
	}()
	for i := len(dynamic[0]) - 2; i >= 0; i-- {
		minIndex, min := 0, math.Inf(1)
		for k := 0; k < (maxStep*2 + 1); k++ {
			tempIndex := (result[i+1] - maxStep + k) % (len(dynamic) - 1)
			if tempIndex < 0 {
				tempIndex = 0
			}
			temp := dynamic[tempIndex][i]
			if temp < min {
				min = temp
				minIndex = tempIndex
			}
		}
		result[i] = minIndex
	}

	return result
}

func GetVerticalDynamicPrepResult(energies [][]float64, maxStep int) [][]float64 {
	// creating 2-dimensional array filled with +Inf
	dynamic := make([][]float64, len(energies))
	for i := 0; i < len(dynamic); i++ {
		dynamic[i] = make([]float64, len(energies[0]))
		for j := 0; j < len(dynamic[i]); j++ {
			dynamic[i][j] = energies[i][j]
		}
	}

	// using dynamic programming to get min cumulative energy for each point starting from left
	for i := 1; i < len(dynamic[0]); i++ {
		for j := 0; j < len(dynamic); j++ {
			min := math.Inf(1)
			for k := 0; k < (maxStep*2 + 1); k++ {
				index := (j - maxStep + k) % (len(dynamic) - 1)
				if index < 0 {
					index = 0
				}
				min = math.Min(dynamic[index][i-1], min)
			}
			dynamic[j][i] = min + energies[j][i]
		}
	}

	return dynamic
}
