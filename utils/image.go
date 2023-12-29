package utils

import (
	"image"
	"image/color"
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
	rows := len(energies)
	cols := len(energies[0])

	dynamic := energies

	for i := 1; i < rows; i++ {
		for j := 0; j < cols; j++ {
			min := math.Inf(1)
			for k := -maxStep; k <= maxStep; k++ {
				index := j + k
				if index < 0 {
					index = 0
				} else if index >= cols {
					index = cols - 1
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
	rows := len(energies)
	cols := len(energies[0])

	dynamic := energies

	for i := 1; i < cols; i++ {
		for j := 0; j < rows; j++ {
			min := math.Inf(1)
			start := j - maxStep
			end := j + maxStep

			if start < 0 {
				start = 0
			}
			if end >= rows {
				end = rows - 1
			}

			for k := start; k <= end; k++ {
				if val := dynamic[k][i-1]; val < min {
					min = val
				}
			}
			dynamic[j][i] = min + energies[j][i]
		}
	}

	return dynamic
}

func CreateImageFromColorMap(colors [][]color.Color, width int, height int) image.Image {
	if width > len(colors) {
		width = len(colors)
	}
	if height > len(colors[0]) {
		height = len(colors[0])
	}
	// Create an RGBA image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Iterate through the colors and set them in the image
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, colors[x][y])
		}
	}

	return img
}
