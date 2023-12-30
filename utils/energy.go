package utils

import (
	"image/color"
	"math"
)

func CalculateE1(img [][]color.Color, x, y int) float64 {
	// Sobel operator for e1 energy calculation along the x-axis
	sobelX := [][]int{{-1, 0, 1}, {-2, 0, 2}, {-1, 0, 1}}

	// Calculate e1 energy for the given pixel (x, y)
	var dx float64

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			// Check boundaries
			if x+i >= 0 && x+i < len(img) && y+j >= 0 && y+j < len(img[x]) {
				// Convert color to grayscale
				gray := color.GrayModel.Convert(img[x+i][y+j]).(color.Gray).Y
				dx += float64(gray) * float64(sobelX[i+1][j+1])
			}
		}
	}

	// Return the absolute value of the calculated e1 energy as float64
	return math.Abs(dx)
}

func CalculateEntropyEnergy(img [][]color.Color, x, y int) float64 {
	// Calculate entropy energy over a 9x9 window centered around the pixel (x, y)
	windowSize := 9
	entropy := 0.0

	for i := -windowSize / 2; i <= windowSize/2; i++ {
		for j := -windowSize / 2; j <= windowSize/2; j++ {
			// Check boundaries
			if x+i >= 0 && x+i < len(img) && y+j >= 0 && y+j < len(img[x]) {
				// Convert color to grayscale
				gray := color.GrayModel.Convert(img[x+i][y+j]).(color.Gray).Y
				entropy -= float64(gray) * math.Log2(float64(gray)/255)
			}
		}
	}

	return entropy * CalculateE1(img, x, y)
}

func gradientMagnitude(dx, dy float64) float64 {
	return math.Sqrt(dx*dx + dy*dy)
}

func CalculateGradient(img [][]color.Color, x, y int) float64 {
	// Sobel operator for gradient calculation
	sobelX := [][]int{{-1, 0, 1}, {-2, 0, 2}, {-1, 0, 1}}
	sobelY := [][]int{{-1, -2, -1}, {0, 0, 0}, {1, 2, 1}}

	var dx, dy float64

	bounds := len(img)
	if boundsX := len(img[0]); boundsX < bounds {
		bounds = boundsX
	}

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			// Check boundaries
			if x+i >= 0 && x+i < bounds && y+j >= 0 && y+j < bounds {
				// Convert color to grayscale
				gray := color.GrayModel.Convert(img[x+i][y+j]).(color.Gray).Y
				dx += float64(gray) * float64(sobelX[i+1][j+1])
				dy += float64(gray) * float64(sobelY[i+1][j+1])
			}
		}
	}

	return gradientMagnitude(dx, dy)
}

func calculateGradient(img [][]color.Color, x, y int) (float64, float64) {
	// Sobel operator for gradient calculation
	sobelX := [][]int{{-1, 0, 1}, {-2, 0, 2}, {-1, 0, 1}}
	sobelY := [][]int{{-1, -2, -1}, {0, 0, 0}, {1, 2, 1}}

	var dx, dy float64

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if x+i >= 0 && x+i < len(img) && y+j >= 0 && y+j < len(img[x]) {
				// Convert color to grayscale
				gray := color.GrayModel.Convert(img[x+i][y+j]).(color.Gray).Y
				dx += float64(gray) * float64(sobelX[i+1][j+1])
				dy += float64(gray) * float64(sobelY[i+1][j+1])
			}
		}
	}

	return math.Abs(dx), math.Abs(dy)
}

func calculateHoG(img [][]color.Color, x, y int) [8]float64 {
	histogram := [8]float64{}

	// Calculate HoG in an 11x11 window centered at (x, y)
	windowSize := 11
	binSize := 256 / 8 // Assuming 8 bins for the histogram

	for i := -windowSize / 2; i <= windowSize/2; i++ {
		for j := -windowSize / 2; j <= windowSize/2; j++ {
			if x+i >= 0 && x+i < len(img) && y+j >= 0 && y+j < len(img[x]) {
				// Convert color to grayscale
				gray := color.GrayModel.Convert(img[x+i][y+j]).(color.Gray).Y

				// Calculate gradient magnitude
				dx, dy := calculateGradient(img, x+i, y+j)
				gradient := gradientMagnitude(dx, dy)

				// Calculate the bin index
				bin := int(gray) / binSize
				if bin >= 8 {
					bin = 7
				}

				histogram[bin] += gradient
			}
		}
	}

	return histogram
}

func CalculateEHoG(img [][]color.Color, x, y int) float64 {
	// Calculate gradients along x and y axes
	dx, dy := calculateGradient(img, x, y)

	// Calculate eHoG
	gradientMax := math.Max(math.Abs(dx), math.Abs(dy))

	// Calculate HoG at the pixel
	histogram := calculateHoG(img, x, y)

	// Find the maximum value in HoG
	maxHoG := histogram[0]
	for _, val := range histogram {
		if val > maxHoG {
			maxHoG = val
		}
	}

	return gradientMax + maxHoG
}
