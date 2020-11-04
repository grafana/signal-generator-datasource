package waves

import (
	"math"
	"math/rand"
	"time"
)

type WaveformArgs struct {
	Type      string    `json:"type,omitempty"`   // in seconds
	Period    float64   `json:"period,omitempty"` // in seconds
	Amplitude float64   `json:"amplitude,omitempty"`
	DutyCycle float64   `json:"duty,omitempty"` // on time vs off time (0-1)
	Points    []float64 `json:"points,omitempty"`
	Ease      string    `json:"ease,omitempty"` // use for animation
}

// Given 0-1 return a scaling function
type WaveformFunc func(t time.Time, args *WaveformArgs) float64

// Registry of known scaling functions
var WaveformFunctions = map[string]WaveformFunc{
	"Sin":      sinFunc,
	"Pulse":    pluseFunc,
	"Sawtooth": sawtoothFunc,
	"Sinc":     sincFunc,
	"Noise":    noiseFunc,
	"CSV":      csvFunc,
}

// Find where in the period things are
func getPeriodPercent(t time.Time, args *WaveformArgs) float64 {
	if args.Period <= 0 {
		return 0
	}
	m := t.UnixNano() % int64(args.Period*1000000000)
	return float64(m) / (args.Period * 1000000000)
}

func sinFunc(t time.Time, args *WaveformArgs) float64 {
	x := getPeriodPercent(t, args)
	return math.Sin(x*2*math.Pi) * args.Amplitude
}

func noiseFunc(t time.Time, args *WaveformArgs) float64 {
	r := rand.New(rand.NewSource(t.UnixNano())) // will be consistent for the value
	return r.Float64() * args.Amplitude
}

func pluseFunc(t time.Time, args *WaveformArgs) float64 {
	p := getPeriodPercent(t, args)
	if p > args.DutyCycle {
		return args.Amplitude
	}
	return 0
}

func sincFunc(t time.Time, args *WaveformArgs) float64 {
	x := getPeriodPercent(t, args)
	return (math.Sin(x) / x) * args.Amplitude
}

func sawtoothFunc(t time.Time, args *WaveformArgs) float64 {
	p := getPeriodPercent(t, args)
	return args.Amplitude * p
}

func csvFunc(t time.Time, args *WaveformArgs) float64 {
	count := float64(len(args.Points))
	if count == 0 {
		return args.Amplitude
	}
	if count <= 1 {
		return args.Amplitude * args.Points[0]
	}

	p := getPeriodPercent(t, args)
	if p >= 1 { // return the last point
		return args.Points[len(args.Points)-1]
	}

	// Step functions to each point
	if args.Ease == "" {
		idx := int(math.Floor(p * count))
		return args.Points[idx]
	}

	f, ok := EaseFunctions[args.Ease]
	if !ok {
		f = EaseLinear
	}

	idx := int(math.Floor(p * (count - 1)))
	step := 1 / (count - 1)
	stepp := p - (step * float64(idx))

	start := args.Points[idx]
	next := args.Points[idx+1]
	delta := next - start

	v := start + (f(stepp) * delta)

	return v * args.Amplitude
}
