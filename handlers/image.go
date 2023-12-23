package handlers

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"net/http"
	L "seam-carving/utils"
)

func Image_handler(w http.ResponseWriter, r *http.Request) {
	img, err := L.GetImageFromRequest(r)
	if err != nil {
		http.Error(w, "Error decoding image: "+err.Error(), http.StatusBadRequest)
		return
	}

	energies := make([][]float64, img.Bounds().Dx())
	for i := range energies {
		energies[i] = make([]float64, img.Bounds().Dy())
	}

	max := 0.0
	for i := 0; i < len(energies); i++ {
		for j := 0; j < len(energies[0]); j++ {
			energy := L.ComputeEnergy(img, i, j)

			energies[i][j] = energy
			if energy > max {
				max = energy
			}
		}
	}

	seam := getHorizontalSeam(energies)

	result := image.NewRGBA(img.Bounds())
	for i := 0; i < len(energies); i++ {
		for j := 0; j < len(energies[0]); j++ {
			if j == seam[i] {
				println(j)
				result.Set(i, j, color.White)
			} else {
				result.Set(i, j, color.Gray16{0xffff - uint16((energies[i][j]/max)*0xffff)})
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	png.Encode(w, result.SubImage(result.Rect))

	// width, err := strconv.Atoi(r.FormValue("width"))
	// if err != nil {
	// 	http.Error(w, "Invalid width value", http.StatusBadRequest)
	// }

	// height, err := strconv.Atoi(r.FormValue("height"))
	// if err != nil {
	// 	http.Error(w, "Invalid height value", http.StatusBadRequest)
	// }

	//fmt.Fprintf(w, "Image size: %d - width, %d - height", img.Bounds().Dx(), img.Bounds().Dy())
}

func getHorizontalSeam(energies [][]float64) []int {
	// creating 2-dimentional array filled with +Inf
	dynamic := make([][]float64, len(energies))
	for i := 0; i < len(dynamic); i++ {
		dynamic[i] = make([]float64, len(energies[0]))
		dynamic[i][0] = energies[i][0]
		for j := 1; j < len(dynamic[i]); j++ {
			dynamic[i][j] = math.Inf(1)
		}
	}

	// using dynamic programming to get min cumulative energy for each point starting from top
	for i := 1; i < len(dynamic); i++ {
		for j := 0; j < len(dynamic[i]); j++ {
			minValue := func() float64 {
				if j == 0 {
					return math.Min(dynamic[i-1][j], dynamic[i-1][j+1])
				} else if j == (len(dynamic[i]) - 1) {
					return math.Min(dynamic[i-1][j-1], dynamic[i-1][j])
				} else {
					return math.Min(dynamic[i-1][j-1], dynamic[i-1][j])
				}
			}()
			dynamic[i][j] = energies[i][j] + minValue
		}
	}

	result := make([]int, len(energies))
	for row := len(energies) - 1; row >= 0; row-- {
		minIndex := 0
		slice := func() []float64 {
			if row > 0 {
				if result[row-1] == 0 {
					return energies[row][result[row-1] : result[row-1]+2]
				} else if result[row-1] == len(energies)-1 {
					return energies[row][result[row-1]-1 : result[row-1]+1]
				} else {
					return energies[row][result[row-1]-1 : result[row-1]+2]
				}
			} else {
				return energies[row]
			}
		}()
		for i, v := range slice {
			if energies[row][minIndex] > v {
				minIndex = i
			}
		}
		result[row] = minIndex
	}

	return result
}
