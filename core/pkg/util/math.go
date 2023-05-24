package util

import "math"

func CalculateDiagonal(s1 float64, s2 float64) float64 {
	return math.Sqrt(math.Pow(s1, 2) + math.Pow(s2, 2))
}
