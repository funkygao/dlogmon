package stats

import (
	"math"
)

// Data structure to contain accumulating values and moments
type Stats struct {
	n, min, max, sum, mean, m2, m3, m4 float64
}

func (d *Stats) Count() int {
	return int(d.n)
}

func (d *Stats) Size() int {
	return int(d.n)
}

func (d *Stats) Min() float64 {
	return d.min
}

func (d *Stats) Max() float64 {
	return d.max
}

func (d *Stats) Sum() float64 {
	return d.sum
}

func (d *Stats) Mean() float64 {
	return d.mean
}

// Update the stats with the given value.
func (d *Stats) Update(x float64) {
	if d.n == 0.0 || x < d.min {
		d.min = x
	}
	if d.n == 0.0 || x > d.max {
		d.max = x
	}
	d.sum += x
	nMinus1 := d.n
	d.n += 1.0
	delta := x - d.mean
	delta_n := delta / d.n
	delta_n2 := delta_n * delta_n
	term1 := delta * delta_n * nMinus1
	d.mean += delta_n
	d.m4 += term1*delta_n2*(d.n*d.n-3*d.n+3.0) + 6*delta_n2*d.m2 - 4*delta_n*d.m3
	d.m3 += term1*delta_n*(d.n-2.0) - 3*delta_n*d.m2
	d.m2 += term1
}

// Update the stats with the given array of values.
func (d *Stats) UpdateArray(data []float64) {
	for _, v := range data {
		d.Update(v)
	}
}

func (d *Stats) PopulationVariance() float64 {
	if d.n == 0 || d.n == 1 {
		return math.NaN()
	}
	return d.m2 / d.n
}

func (d *Stats) SampleVariance() float64 {
	if d.n == 0 || d.n == 1 {
		return math.NaN()
	}
	return d.m2 / (d.n - 1.0)
}

func (d *Stats) PopulationStandardDeviation() float64 {
	if d.n == 0 || d.n == 1 {
		return math.NaN()
	}
	return math.Sqrt(d.PopulationVariance())
}

func (d *Stats) SampleStandardDeviation() float64 {
	if d.n == 0 || d.n == 1 {
		return math.NaN()
	}
	return math.Sqrt(d.SampleVariance())
}

func (d *Stats) PopulationSkew() float64 {
	return math.Sqrt(d.n/(d.m2*d.m2*d.m2)) * d.m3
}

func (d *Stats) SampleSkew() float64 {
	if d.n == 2.0 {
		return math.NaN()
	}
	popSkew := d.PopulationSkew()
	return math.Sqrt(d.n*(d.n-1.0)) / (d.n - 2.0) * popSkew
}

// The kurtosis functions return _excess_ kurtosis, so that the kurtosis of a normal
// distribution = 0.0. Then kurtosis < 0.0 indicates platykurtic (flat) while
// kurtosis > 0.0 indicates leptokurtic (peaked) and near 0 indicates mesokurtic.Update
func (d *Stats) PopulationKurtosis() float64 {
	return (d.n*d.m4)/(d.m2*d.m2) - 3.0
}

func (d *Stats) SampleKurtosis() float64 {
	if d.n == 2.0 || d.n == 3.0 {
		return math.NaN()
	}
	populationKurtosis := d.PopulationKurtosis()
	return (d.n - 1.0) / ((d.n - 2.0) * (d.n - 3.0)) * ((d.n+1.0)*populationKurtosis + 6.0)
}
