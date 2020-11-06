package waves

import (
	"math"
	"math/rand"
	"time"
)

type WaveformArgs struct {
	Type      string    `json:"type,omitempty"`
	Period    string    `json:"period,omitempty"`    // parse duration or range/X
	PeriodSec float64   `json:"periodSec,omitempty"` // in seconds
	Amplitude float64   `json:"amplitude,omitempty"`
	Offset    float64   `json:"offset,omitempty"`
	Phase     float64   `json:"phase,omitempty"`
	DutyCycle float64   `json:"duty,omitempty"` // on time vs off time (0-1)
	Points    []float64 `json:"points,omitempty"`
	Args      string    `json:"args,omitempty"` // ease function or expression
}

// Given 0-1 return a scaling function -- note this does not include amplitude and doffset
type WaveformFunc func(t time.Time, args *WaveformArgs) float64

// Registry of known scaling functions
var WaveformFunctions = map[string]WaveformFunc{
	"Sin":      sinFunc,
	"Square":   squareFunc,
	"Triangle": triangleFunc,
	"Sawtooth": sawtoothFunc,
	"Noise":    noiseFunc,
	"CSV":      csvFunc,
}

// Find where in the period things are
func getPeriodPercent(t time.Time, args *WaveformArgs) float64 {
	if args.PeriodSec <= 0 {
		return 0
	}
	m := t.UnixNano() % int64(args.PeriodSec*1000000000)
	p := (float64(m) / (args.PeriodSec * 1000000000)) + args.Phase
	if p > 1 {
		return p - 1 // wrap the phase
	}
	return p
}

func sinFunc(t time.Time, args *WaveformArgs) float64 {
	x := getPeriodPercent(t, args)
	return math.Sin(x * 2 * math.Pi)
}

func noiseFunc(t time.Time, args *WaveformArgs) float64 {
	r := rand.New(rand.NewSource(t.UnixNano())) // will be consistent for the value
	return (r.Float64() * 2) - 1
}

func squareFunc(t time.Time, args *WaveformArgs) float64 {
	p := getPeriodPercent(t, args)
	if p > args.DutyCycle {
		return 1
	}
	return -1
}

func triangleFunc(t time.Time, args *WaveformArgs) float64 {
	p := getPeriodPercent(t, args)
	if p > 0.75 {
		return ((p - .75) * 4) - 1
	}
	if p > 0.25 {
		return 1 - ((p - .25) * 4)
	}
	return (p * 4)
}

func sawtoothFunc(t time.Time, args *WaveformArgs) float64 {
	p := getPeriodPercent(t, args)
	return (p * 2) - 1
}

func csvFunc(t time.Time, args *WaveformArgs) float64 {
	count := float64(len(args.Points))
	if count == 0 {
		return 0 // center at zero
	}
	if count <= 1 {
		return args.Points[0]
	}

	p := getPeriodPercent(t, args)
	if p >= 1 { // return the last point
		return args.Points[len(args.Points)-1]
	}

	// Step functions to each point
	if args.Args == "" {
		idx := int(math.Floor(p * count))
		return args.Points[idx]
	}

	f, ok := EaseFunctions[args.Args]
	if !ok {
		f = EaseLinear
	}

	idx := int(math.Floor(p * (count - 1)))
	step := 1 / (count - 1)
	stepp := p - (step * float64(idx))

	start := args.Points[idx]
	next := args.Points[idx+1]
	delta := next - start

	return start + (f(stepp) * delta)
}
