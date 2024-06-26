package Canvas

import (
	"fmt"
	"math"
)

const (
	AnsiMode_256 = iota
	AnsiMode_TrueColor
)

type RGB struct {
	R byte
	G byte
	B byte
}

func (rgb *RGB) to256Code() int {
	// The six intensity levels
	levels := []int{0, 95, 135, 175, 215, 255}

	// Function to find the closest level
	closest := func(value uint8) int {
		closestLevel := levels[0]
		minDistance := math.Abs(float64(value) - float64(levels[0]))
		for _, level := range levels {
			distance := math.Abs(float64(value) - float64(level))
			if distance < minDistance {
				minDistance = distance
				closestLevel = level
			}
		}
		return closestLevel
	}

	// Finding the closest levels for r, g, b
	redLevel := closest(rgb.R)
	greenLevel := closest(rgb.G)
	blueLevel := closest(rgb.B)

	// Map levels to their index in the 0-5 range
	mapToIndex := func(level int) int {
		for i, l := range levels {
			if l == level {
				return i
			}
		}
		return 0 // default, should not reach here
	}

	redIndex := mapToIndex(redLevel)
	greenIndex := mapToIndex(greenLevel)
	blueIndex := mapToIndex(blueLevel)

	// Calculating the final index
	colourCubeIndex := 16 + 36*redIndex + 6*greenIndex + blueIndex
	return colourCubeIndex
}

func (rgb *RGB) To256AnsiFG() string {
	code := rgb.to256Code()
	return fmt.Sprintf("\033[38;5;%dm", code)
}
func (rgb *RGB) To256AnsiBG() string {
	code := rgb.to256Code()
	return fmt.Sprintf("\033[48;5;%dm", code)
}
func (rgb *RGB) ToTrueColorFG() string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", rgb.R, rgb.G, rgb.B)
}
func (rgb *RGB) ToTrueColorBG() string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", rgb.R, rgb.G, rgb.B)
}
