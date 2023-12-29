package utils

import (
	"image/color"
	"math"
)

func ComputeEnergy(colorMap [][]color.Color, x int, y int) float64 {
	width := len(colorMap)
	height := len(colorMap[0])

	rx0, gx0, bx0, _ := colorMap[func() int {
		if x == 0 {
			return width - 1
		}
		return x - 1
	}()][y].RGBA()
	rx1, gx1, bx1, _ := colorMap[func() int {
		if x == width-1 {
			return 0
		}
		return x + 1
	}()][y].RGBA()

	deltaX := math.Pow((float64)(rx1-rx0), 2) + math.Pow((float64)(gx1-gx0), 2) + math.Pow((float64)(bx1-bx0), 2)

	ry0, gy0, by0, _ := colorMap[x][func() int {
		if y == 0 {
			return height - 1
		}
		return y - 1
	}()].RGBA()
	ry1, gy1, by1, _ := colorMap[x][func() int {
		if y == height-1 {
			return 0
		}
		return y + 1
	}()].RGBA()

	deltaY := math.Pow((float64)(ry1-ry0), 2) + math.Pow((float64)(gy1-gy0), 2) + math.Pow((float64)(by1-by0), 2)

	return deltaX + deltaY
}
