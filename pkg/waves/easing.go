package waves

import "math"

const c1 = 1.70158
const c2 = c1 * 1.525
const c3 = c1 + 1
const c4 = (2 * math.Pi) / 3
const c5 = (2 * math.Pi) / 4.5

// Given 0-1 return a scaling function
type EaseFunc func(x float64) float64

// Registry of known scaling functions
var EaseFunctions = map[string]EaseFunc{
	"Linear":       EaseLinear,
	"InQuad":       easeInQuad,
	"OutQuad":      easeOutQuad,
	"InOutQuad":    easeInOutQuad,
	"InCubic":      easeInCubic,
	"OutCubic":     easeOutCubic,
	"InOutCubic":   easeInOutCubic,
	"InQuart":      easeInQuart,
	"OutQuart":     easeOutQuart,
	"InOutQuart":   easeInOutQuart,
	"InQuint":      easeInQuint,
	"OutQuint":     easeOutQuint,
	"InOutQuint":   easeInOutQuint,
	"InSine":       easeInSine,
	"OutSine":      easeOutSine,
	"InOutSine":    easeInOutSine,
	"InExpo":       easeInExpo,
	"OutExpo":      easeOutExpo,
	"InOutExpo":    easeInOutExpo,
	"InCirc":       easeInCirc,
	"OutCirc":      easeOutCirc,
	"InOutCirc":    easeInOutCirc,
	"InBack":       easeInBack,
	"OutBack":      easeOutBack,
	"InOutBack":    easeInOutBack,
	"InElastic":    easeInElastic,
	"OutElastic":   easeOutElastic,
	"InOutElastic": easeInOutElastic,
}

// EaseLinear simple linear easing
func EaseLinear(x float64) float64 {
	return x
}

func easeInQuad(x float64) float64 {
	return x * x
}

func easeOutQuad(x float64) float64 {
	return 1 - (1-x)*(1-x)
}

func easeInOutQuad(x float64) float64 {
	if x < 0.5 {
		return 2 * x * x
	}
	return 1 - math.Pow(-2*x+2, 2)/2.0
}

func easeInCubic(x float64) float64 {
	return x * x * x
}

func easeOutCubic(x float64) float64 {
	return 1 - math.Pow(1-x, 3)
}

func easeInOutCubic(x float64) float64 {
	if x < 0.5 {
		return 4 * x * x * x
	}
	return 1 - math.Pow(-2*x+2, 3)/2
}

func easeInQuart(x float64) float64 {
	return x * x * x * x
}

func easeOutQuart(x float64) float64 {
	return 1 - math.Pow(1-x, 4)
}

func easeInOutQuart(x float64) float64 {
	if x < 0.5 {
		return 8 * x * x * x * x
	}
	return 1 - math.Pow(-2*x+2, 4)/2
}

func easeInQuint(x float64) float64 {
	return x * x * x * x * x
}

func easeOutQuint(x float64) float64 {
	return 1 - math.Pow(1-x, 5)
}

func easeInOutQuint(x float64) float64 {
	if x < 0.5 {
		return 16 * x * x * x * x * x
	}
	return 1 - math.Pow(-2*x+2, 5)/2
}

func easeInSine(x float64) float64 {
	return 1 - math.Cos((x*math.Pi)/2)
}

func easeOutSine(x float64) float64 {
	return math.Sin((x * math.Pi) / 2)
}

func easeInOutSine(x float64) float64 {
	return -(math.Cos(math.Pi*x) - 1) / 2
}

func easeInExpo(x float64) float64 {
	if x == 0 {
		return 0
	}
	return math.Pow(2, 10*x-10)
}

func easeOutExpo(x float64) float64 {
	if x == 1 {
		return 1
	}
	return 1 - math.Pow(2, -10*x)
}

func easeInOutExpo(x float64) float64 {
	if x == 0 {
		return 0
	}
	if x == 1 {
		return 1
	}
	if x < 0.5 {
		return math.Pow(2, 20*x-10) / 2
	}
	return (2 - math.Pow(2, -20*x+10)) / 2
}

func easeInCirc(x float64) float64 {
	return 1 - math.Sqrt(1-math.Pow(x, 2))
}

func easeOutCirc(x float64) float64 {
	return math.Sqrt(1 - math.Pow(x-1, 2))
}

func easeInOutCirc(x float64) float64 {
	if x < 0.5 {
		return (1 - math.Sqrt(1-math.Pow(2*x, 2))) / 2
	}
	return (math.Sqrt(1-math.Pow(-2*x+2, 2)) + 1) / 2
}

func easeInBack(x float64) float64 {
	return c3*x*x*x - c1*x*x
}

func easeOutBack(x float64) float64 {
	return 1 + c3*math.Pow(x-1, 3) + c1*math.Pow(x-1, 2)
}

func easeInOutBack(x float64) float64 {
	if x < 0.5 {
		return (math.Pow(2*x, 2) * ((c2+1)*2*x - c2)) / 2
	}
	return (math.Pow(2*x-2, 2)*((c2+1)*(x*2-2)+c2) + 2) / 2
}

func easeInElastic(x float64) float64 {
	if x == 0 {
		return 0
	}
	if x == 1 {
		return 1
	}
	return -math.Pow(2, 10*x-10) * math.Sin((x*10-10.75)*c4)
}

func easeOutElastic(x float64) float64 {
	if x == 0 {
		return 0
	}
	if x == 1 {
		return 1
	}
	return math.Pow(2, -10*x)*math.Sin((x*10-0.75)*c4) + 1
}

func easeInOutElastic(x float64) float64 {
	if x == 0 {
		return 0
	}
	if x == 1 {
		return 1
	}
	if x < 0.5 {
		return -(math.Pow(2, 20*x-10) * math.Sin((20*x-11.125)*c5)) / 2
	}
	return (math.Pow(2, -20*x+10)*math.Sin((20*x-11.125)*c5))/2 + 1
}
