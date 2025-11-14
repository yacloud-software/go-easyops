package utils

import "testing"

func TestInterpolate(t *testing.T) {
	testinterpol(t, map[float64]float64{3: 30, 5: 50, 10: 100}, 7, 70)
	testinterpol(t, map[float64]float64{3: 30, 5: 50, 10: 100}, 8, 80)
	testinterpol(t, map[float64]float64{3: 30, 5: 50, 10: 100}, 4, 40)
	testinterpol(t, map[float64]float64{3: 30, 5: 50, 10: 100}, 2, 20)
	testinterpol(t, map[float64]float64{3: 30, 5: 50, 10: 100}, 11, 100)

	testinterpol(t, map[float64]float64{3: 30, 5: 50, 7: 100}, 4, 40)
	testinterpol(t, map[float64]float64{3: 30, 5: 50, 7: 100}, 6, 75)

}
func testinterpol(t *testing.T, ipm map[float64]float64, num, expected float64) {
	ip := &Interpolator{}
	ip.AddReferencePoints(ipm)
	res := ip.LinearInterpolate(num)
	if res == expected {
		return
	}
	t.Logf("For interpolator (%s), value %0.1f, expected %0.1f, but got %0.1f", ip.String(), num, expected, res)
	t.Fail()

}
