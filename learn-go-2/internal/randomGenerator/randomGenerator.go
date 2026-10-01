package randomgenerator

import "math/rand"

func Random1(minimalValue int, maxValue int) int {
	return rand.Int()%
		(maxValue-minimalValue) +
		minimalValue
}

func Random2(minimalValue int, maxValue int) int {
	return int(
		rand.Float32()*
			float32(maxValue-minimalValue),
	) + minimalValue
}
