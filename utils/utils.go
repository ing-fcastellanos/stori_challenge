package utils

import "strconv"

func ParseInt(str string) int {
	val, _ := strconv.Atoi(str)
	return val
}

func ParseFloat(str string) float64 {
	val, _ := strconv.ParseFloat(str, 64)
	return val
}
