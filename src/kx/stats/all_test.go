package stats

import (
	"github.com/bmizerany/assert"
	"testing"
)

func mockStats() Stats {
	return Stats{}
}

func TestStatsNormalUsage(t *testing.T) {
	v := []float64{4, 19, 904}
	s := mockStats()
	s.UpdateArray(v)
	assert.Equal(t, 927.0, s.Sum())
	assert.Equal(t, 309.0, s.Mean())
	assert.Equal(t, 904.0, s.Max())
	assert.Equal(t, 4.0, s.Min())
	assert.Equal(t, 3, s.Count())
	assert.Equal(t, 515.3396937942972, s.SampleStandardDeviation())
}

func TestFuncUsage(t *testing.T) {
	v := []float64{4, 19, 904}
	assert.Equal(t, 927.0, StatsSum(v))
	assert.Equal(t, 309.0, StatsMean(v))
	assert.Equal(t, 904.0, StatsMax(v))
	assert.Equal(t, 4.0, StatsMin(v))
	assert.Equal(t, 3, StatsCount(v))
	assert.Equal(t, 515.3396937942972, StatsSampleStandardDeviation(v))
}
