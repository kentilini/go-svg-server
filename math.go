package main

//deal with empty arrays
func MinMax(array []float32) (float32, float32) {
	var max float32 = array[0]
	var min float32 = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}
