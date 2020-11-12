package waves

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/Knetic/govaluate"
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
var WaveformFunctions = map[string]govaluate.ExpressionFunction{
	"Sine":     sineFunc,
	"Square":   squareFunc,
	"Triangle": triangleFunc,
	"Sawtooth": sawtoothFunc,
	"Noise":    noiseFunc,
	"CSV":      csvFunc,
}

func sineFunc(args ...interface{}) (interface{}, error) {
	x, ok := args[0].(float64)
	if !ok {
		return 0, fmt.Errorf("should be a number!")
	}
	return math.Sin(x), nil
}

var grand = rand.New(rand.NewSource(0))

func noiseFunc(args ...interface{}) (interface{}, error) {
	r := grand
	if len(args) > 0 {
		x, ok := args[0].(int64)
		if !ok {
			return 0, fmt.Errorf("invalid seed")
		}
		r = rand.New(rand.NewSource(x)) // will be consistent for the value
	}
	return (r.Float64() * 2) - 1, nil
}

func squareFunc(args ...interface{}) (interface{}, error) {
	x, ok := args[0].(float64)
	if !ok {
		return 0, fmt.Errorf("missing value")
	}
	dutyCycle := 0.5
	if len(args) > 1 {
		dutyCycle, ok = args[1].(float64)
		if !ok {
			return 0, fmt.Errorf("missing 0-1 duty")
		}
		if dutyCycle > 1 {
			dutyCycle = 1
		} else if dutyCycle < 0 {
			dutyCycle = 0
		}
	}
	p := (x) / (2 * math.Pi)
	if p > dutyCycle {
		return -1.0, nil
	}
	return 1.0, nil
}

func triangleFunc(args ...interface{}) (interface{}, error) {
	x, ok := args[0].(float64)
	if !ok {
		return 0, fmt.Errorf("missing value")
	}

	p := (x) / (2.0 * math.Pi)
	if p > 0.75 {
		return ((p - .75) * 4.0) - 1.0, nil
	}
	if p > 0.25 {
		return 1 - ((p - .25) * 4.0), nil
	}
	return (p * 4.0), nil

	// p := (x + math.Pi) / (2 * math.Pi)
	// if p > 0.75 {
	// 	return 4 - (p * 4.0), nil
	// }
	// if p > 0.25 {
	// 	return ((p - .25) * 4) - 1.0, nil // GOOD
	// }
	// return (-p * 4.0), nil
}

func sawtoothFunc(args ...interface{}) (interface{}, error) {
	x, ok := args[0].(float64)
	if !ok {
		return 0, fmt.Errorf("missing value")
	}
	return (x / math.Pi) - 1, nil
}

func csvFunc(args ...interface{}) (interface{}, error) {

	return 0, nil

	// count := float64(len(args.Points))
	// if count == 0 {
	// 	return 0 // center at zero
	// }
	// if count <= 1 {
	// 	return args.Points[0]
	// }

	// p := getPeriodPercent(t, args)
	// if p >= 1 { // return the last point
	// 	return args.Points[len(args.Points)-1]
	// }

	// // Step functions to each point
	// if args.Args == "" {
	// 	idx := int(math.Floor(p * count))
	// 	return args.Points[idx]
	// }

	// f, ok := EaseFunctions[args.Args]
	// if !ok {
	// 	f = EaseLinear
	// }

	// idx := int(math.Floor(p * (count - 1)))
	// step := 1 / (count - 1)
	// stepp := p - (step * float64(idx))

	// start := args.Points[idx]
	// next := args.Points[idx+1]
	// delta := next - start

	// return start + (f(stepp) * delta)
}
