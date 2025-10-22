package utils

import (
	"fmt"
	"strconv"
)

func ParseInt(str string) (int, error) {
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("valor no valido")
	}
	return val, nil
}

func ParseFloat(str string) (float64, error) {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, fmt.Errorf("valor no valido")
	}
	return val, nil
}
