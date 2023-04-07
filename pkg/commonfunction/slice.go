package commonfunction

import "math"

func CommonElement(list1, list2 []string) []string {

	var matches []string

	for _, item1 := range list1 {
		for _, item2 := range list2 {
			if item1 == item2 {
				matches = append(matches, item1)
			}
		}
	}

	return matches

}

func WeightedStandardDeviation(data []float64, weights []float64) float64 {
	n := len(data)
	if n == 0 || n != len(weights) {
		return math.NaN()
	}

	var sum float64
	var sumWeights float64
	for i := 0; i < n; i++ {
		sum += weights[i] * data[i]
		sumWeights += weights[i]
	}
	mean := sum / sumWeights

	var sumSquaredDeviations float64
	for i := 0; i < n; i++ {
		deviation := data[i] - mean
		squaredDeviation := deviation * deviation
		weightedSquaredDeviation := weights[i] * squaredDeviation
		sumSquaredDeviations += weightedSquaredDeviation
	}

	return math.Sqrt(sumSquaredDeviations / sumWeights)
}
