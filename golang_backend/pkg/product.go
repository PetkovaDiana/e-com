package pkg

import (
	"math"
	"reflect"
	"strconv"
	"strings"
)

func ToArray(str string, tp reflect.Type) interface{} {
	chunks := strings.Split(str, ",")

	switch tp.Kind() {
	case reflect.Int:
		res := make([]int, len(chunks))
		for i, c := range chunks {
			res[i], _ = strconv.Atoi(c) // error handling omitted for concision
		}
		return res
	case reflect.String:
		res := make([]string, len(chunks))
		for i, c := range chunks {
			res[i] = c
		}
		return res
	default:
		return nil
	}
}

func Round(x float64) float64 {
	var rounder float64
	pow := math.Pow(10, float64(2))
	intermed := x * pow
	_, frac := math.Modf(intermed)
	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}
	return rounder / pow
}
