package waves

import (
	"testing"
)

// func getValues(check []time.Time, f WaveformFunc, args *WaveformArgs) []float64 {
// 	res := make([]float64, len(check))
// 	for idx, t := range check {
// 		res[idx] = f(t, args)
// 	}
// 	return res
// }

func TestCSV(t *testing.T) {
	// 	args := &WaveformArgs{
	// 		PeriodSec: 10, // 10 seconds
	// 		Amplitude: 1,
	// 	}
	// 	start := time.Unix(0, 0)

	// 	// zero points return amplitude
	// 	v := csvFunc(start, args)
	// 	if v != 0 {
	// 		t.Fail()
	// 	}

	// 	// one point is a constant
	// 	args.Points = []float64{10}
	// 	v = csvFunc(start, args)
	// 	if v != 10 {
	// 		t.Fail()
	// 	}

	// 	// no ease should be a step function
	// 	args.Points = []float64{10, 20}
	// 	if csvFunc(start, args) != 10 {
	// 		t.Fail()
	// 	}
	// 	if csvFunc(start.Add(time.Second*5), args) != 20 {
	// 		t.Fail()
	// 	}
	// 	if csvFunc(start.Add(time.Millisecond*9999), args) != 20 { // almost equal
	// 		t.Fail()
	// 	}

	// 	check := []time.Time{
	// 		start,
	// 		start.Add(time.Second * 4),
	// 		start.Add(time.Second * 5),
	// 		start.Add(time.Millisecond * 9999),
	// 		start.Add(time.Second * 10), // goes back to 0
	// 	}
	// 	vals := getValues(check, csvFunc, args)
	// 	expect := []float64{
	// 		10, 10, 20, 20, 10,
	// 	}
	// 	if diff := cmp.Diff(expect, vals); diff != "" {
	// 		t.Fatalf("unexpect results (-want +got):\n%s", diff)
	// 	}

	// 	// Now with easing
	// 	args.Args = "Linear"
	// 	v = csvFunc(start.Add(time.Second*5), args)
	// 	if diff := cmp.Diff(float64(15), v); diff != "" {
	// 		t.Fatalf("unexpect results (-want +got):\n%s", diff)
	// 	}

	// 	vals = getValues(check, csvFunc, args)
	// 	expect = []float64{
	// 		10, 14, 15, 19.999000000000002, 10,
	// 	}
	// 	if diff := cmp.Diff(expect, vals); diff != "" {
	// 		t.Errorf("unexpect results (-want +got):\n%s", diff)
	// 	}

	// 	// More points
	// 	args.Points = []float64{10, 15, 20}
	// 	v = csvFunc(start.Add(time.Second*5), args)
	// 	if diff := cmp.Diff(float64(15), v); diff != "" {
	// 		t.Fatalf("unexpect results (-want +got):\n%s", diff)
	// 	}
	// 	vals = getValues(check, csvFunc, args)
	// 	expect = []float64{
	// 		10, 12, 15, 17.4995, 10,
	// 	}
	// 	if diff := cmp.Diff(expect, vals); diff != "" {
	// 		t.Errorf("unexpect results (-want +got):\n%s", diff)
	// 	}
	// }

	// func TestNoise(t *testing.T) {
	// 	args := &WaveformArgs{
	// 		PeriodSec: 10, // 10 seconds
	// 		Amplitude: 1,
	// 	}
	// 	start := time.Unix(12345, 0)
	// 	check := []time.Time{
	// 		start,
	// 		start.Add(time.Second * 4),
	// 		start.Add(time.Second * 5),
	// 		start.Add(time.Millisecond * 9999),
	// 		start.Add(time.Second * 10),
	// 	}
	// 	vals := getValues(check, noiseFunc, args)
	// 	// Although random -- they should not change!
	// 	expect := []float64{
	// 		0.07349823459147742, -0.6467401322395636, -0.4666797023040796, -0.006100840921510309, -0.9311339452529835,
	// 	}
	// 	if diff := cmp.Diff(expect, vals); diff != "" {
	// 		t.Fatalf("unexpect results (-want +got):\n%s", diff)
	// 	}
}
