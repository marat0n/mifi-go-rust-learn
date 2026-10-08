package services

import "math/rand"

func Random(minimalValue int, maxValue int) int {
	return rand.Int()%
		(maxValue-minimalValue) +
		minimalValue
}
