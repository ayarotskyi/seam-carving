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

func GetHorizontalSeam(acummulativeMatrix [][]float64, maxStep int) []int {
	width := len(acummulativeMatrix)
	height := len(acummulativeMatrix[0])
	result := make([]int, width)
	result[width-1] = func() int {
		min, minIndex := math.Inf(1), 0
		for i := 0; i < height; i++ {
			temp := acummulativeMatrix[width-1][i]
			if temp < min {
				min = temp
				minIndex = i
			}
		}
		return minIndex
	}()
	for i := len(acummulativeMatrix) - 2; i >= 0; i-- {
		minIndex, min := 0, math.Inf(1)
		for k := 0; k < (maxStep*2 + 1); k++ {

			tempIndex := (result[i+1] - maxStep + k)
			if tempIndex < 0 {
				tempIndex = 0
			}
			if tempIndex >= height {
				tempIndex = height - 1
			}
			temp := acummulativeMatrix[i][tempIndex]
			if temp < min {
				min = temp
				minIndex = tempIndex
			}
		}
		result[i] = minIndex
	}

	return result
}

func GetHorizontalAcummulativeMatrix(energies [][]float64, maxStep int) [][]float64 {
	rows := len(energies)
	cols := len(energies[0])

	acummulativeMatrix := make([][]float64, rows)

	for i := range acummulativeMatrix {
		acummulativeMatrix[i] = make([]float64, cols)
	}
	for j := 0; j < cols; j++ {
		acummulativeMatrix[0][j] = energies[0][j]
	}

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
				nextVal := acummulativeMatrix[i-1][index]
				if nextVal < min {
					min = nextVal
				}
			}
			acummulativeMatrix[i][j] = min + energies[i][j]
		}
	}

	return acummulativeMatrix
}

func GetVerticalSeam(acummulativeMatrix [][]float64, maxStep int) []int {
	width := len(acummulativeMatrix)
	height := len(acummulativeMatrix[0])
	result := make([]int, height)
	result[height-1] = func() int {
		min, minIndex := math.Inf(1), 0
		for i := 0; i < width; i++ {
			temp := acummulativeMatrix[i][height-1]
			if temp < min {
				min = temp
				minIndex = i
			}
		}
		return minIndex
	}()
	for i := height - 2; i >= 0; i-- {
		minIndex, min := 0, math.Inf(1)
		for k := 0; k < (maxStep*2 + 1); k++ {
			tempIndex := (result[i+1] - maxStep + k)
			if tempIndex < 0 {
				tempIndex = 0
			}
			if tempIndex >= width {
				tempIndex = width - 1
			}
			temp := acummulativeMatrix[tempIndex][i]
			if temp < min {
				min = temp
				minIndex = tempIndex
			}
		}
		result[i] = minIndex
	}

	return result
}

func GetVerticalAcummulativeMatrix(energies [][]float64, maxStep int) [][]float64 {
	rows := len(energies)
	cols := len(energies[0])

	acummulativeMatrix := make([][]float64, rows)

	for i := range acummulativeMatrix {
		acummulativeMatrix[i] = make([]float64, cols)
		acummulativeMatrix[i][0] = energies[i][0]
	}

	for i := 0; i < rows; i++ {
		for j := 1; j < cols; j++ {
			min := math.Inf(1)
			for k := -maxStep; k <= maxStep; k++ {
				index := i + k
				if index < 0 {
					index = 0
				} else if index >= rows {
					index = rows - 1
				}
				nextVal := acummulativeMatrix[index][j-1]
				if nextVal < min {
					min = nextVal
				}
			}
			acummulativeMatrix[i][j] = min + energies[i][j]
		}
	}

	return acummulativeMatrix
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
