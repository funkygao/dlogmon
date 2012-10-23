package stats

import (
    "math"
)

func StatsCount(data []float64) int {
	return len(data)
}

func StatsMin(data []float64) float64 {
	if len(data) == 0 {
		return math.NaN()
	}
	min := data[0]
	for _, v := range data {
		if v < min {
			min = v
		}
	}
	return min
}

func StatsMax(data []float64) float64 {
	if len(data) == 0 {
		return math.NaN()
	}
	max := data[0]
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	return max
}

func StatsSum(data []float64) (sum float64) {
	for _, v := range data {
		sum += v
	}
	return
}

func StatsMean(data []float64) float64 {
	return StatsSum(data) / float64(len(data))
}

func sumSquaredDeltas(data []float64) (ssd float64) {
	mean := StatsMean(data)
	for _, v := range data {
		delta := v - mean
		ssd += delta * delta
	}
	return
}

func StatsPopulationVariance(data []float64) float64 {
	n := float64(len(data))
	ssd := sumSquaredDeltas(data)
	return ssd / n
}

func StatsSampleVariance(data []float64) float64 {
	n := float64(len(data))
	ssd := sumSquaredDeltas(data)
	return ssd / (n - 1.0)
}

func StatsPopulationStandardDeviation(data []float64) float64 {
	return math.Sqrt(StatsPopulationVariance(data))
}

func StatsSampleStandardDeviation(data []float64) float64 {
	return math.Sqrt(StatsSampleVariance(data))
}

func StatsPopulationSkew(data []float64) (skew float64) {
	mean := StatsMean(data)
	n := float64(len(data))

	sum3 := 0.0
	for _, v := range data {
		delta := v - mean
		sum3 += delta * delta * delta
	}

	variance := math.Sqrt(StatsPopulationVariance(data))
	skew = sum3 / n / (variance * variance * variance)
	return
}

func StatsSampleSkew(data []float64) float64 {
	popSkew := StatsPopulationSkew(data)
	n := float64(len(data))
	return math.Sqrt(n*(n-1.0)) / (n - 2.0) * popSkew
}

// The kurtosis functions return _excess_ kurtosis
func StatsPopulationKurtosis(data []float64) (kurtosis float64) {
	mean := StatsMean(data)
	n := float64(len(data))

	sum4 := 0.0
	for _, v := range data {
		delta := v - mean
		sum4 += delta * delta * delta * delta
	}

	variance := StatsPopulationVariance(data)
	kurtosis = sum4/(variance*variance)/n - 3.0
	return
}

func StatsSampleKurtosis(data []float64) float64 {
	populationKurtosis := StatsPopulationKurtosis(data)
	n := float64(len(data))
	return (n - 1.0) / ((n - 2.0) * (n - 3.0)) * ((n+1.0)*populationKurtosis + 6.0)
}
