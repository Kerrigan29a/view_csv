package main

import (
	"math"
)

func minSliceFloat64(values []float64) float64 {
	value := values[0]
	nValues := len(values)
	for i := 1; i < nValues; i++ {
		value = math.Min(value, values[i])
	}
	return value
}

func maxSliceFloat64(values []float64) float64 {
	value := values[0]
	nValues := len(values)
	for i := 1; i < nValues; i++ {
		value = math.Max(value, values[i])
	}
	return value
}

func minSliceInt(values []int) int {
	value := values[0]
	nValues := len(values)
	for i := 1; i < nValues; i++ {
		value = minInt(value, values[i])
	}
	return value
}

func maxSliceInt(values []int) int {
	value := values[0]
	nValues := len(values)
	for i := 1; i < nValues; i++ {
		value = maxInt(value, values[i])
	}
	return value
}

func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func normalizeFloat64(values []float64) []float64 {
	nValues := len(values)
	minValue := minSliceFloat64(values)
	maxValue := maxSliceFloat64(values)
	base := maxValue - minValue
	normalizedValues := make([]float64, nValues)
	for i := 0; i < nValues; i++ {
		normalizedValues[i] = (values[i] - minValue) / base
	}
	return normalizedValues
}

func floatToInt(values []float64, base int) []int {
	nValues := len(values)
	intValues := make([]int, nValues)
	for i := 0; i < nValues; i++ {
		intValues[i] = int(math.Round(values[i] * float64(base)))
	}
	return intValues
}

func rescaleInt(values []int) []int {
	nValues := len(values)
	minValue := minSliceInt(values)
	rescaledValues := make([]int, nValues)
	for i := 0; i < nValues; i++ {
		rescaledValues[i] = (values[i] - minValue)
	}
	return rescaledValues
}
