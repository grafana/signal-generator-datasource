import { SelectableValue } from '@grafana/data';
import { WaveformType } from 'types';

export const waveformTypes: Array<SelectableValue<WaveformType>> = [
  {
    label: 'Sine Wave',
    value: WaveformType.Sin,
    description: 'A curve that describes a smooth periodic oscillation',
  },
  {
    label: 'Square Wave',
    value: WaveformType.Square,
    description: 'fixed between minimum and maximum values',
  },
  {
    label: 'Triangle Wave',
    value: WaveformType.Triangle,
    description: 'fixed between minimum and maximum values',
  },
  {
    label: 'Sawtooth Wave',
    value: WaveformType.Sawtooth,
    description: 'fixed between minimum and maximum values',
  },
  {
    label: 'sin(x)/x',
    value: WaveformType.Sawtooth,
    description: 'periodic spikes',
  },
  {
    label: 'Random noise',
    value: WaveformType.Noise,
    description: 'random values',
  },
  {
    label: 'CSV Values',
    value: WaveformType.CSV,
    description: 'Animated values',
  },
];

// "":       EaseLinear,
// "":       ,
// "":      easeOutQuad,
// "":    easeInOutQuad,
// "":      easeInCubic,
// "":     easeOutCubic,
// "":   easeInOutCubic,
// "":      easeInQuart,
// "":     easeOut,
// "":   easeInOutQuart,
// "":      easeIn,
// "":     easeOutQuint,
// "":   easeInOutQuint,
// "":       easeInSine,
// "":      easeOutSine,
// "":    easeInOut,
// "":       easeInExpo,
// "":      easeOutExpo,
// "":    easeInOut,
// "":       easeInCirc,
// "":      easeOutCirc,
// "":    easeInOut,
// "":       easeIn,
// "":      easeOutBack,
// "":    easeInOutBack,

export const easeFunctions: Array<SelectableValue<string>> = [
  {
    label: 'Linear',
    value: 'Linear',
  },
  {
    label: 'InQuad',
    value: 'InQuad',
  },
  {
    label: 'OutQuad',
    value: 'OutQuad',
  },
  {
    label: 'InOutQuad',
    value: 'InOutQuad',
  },
  {
    label: 'InQuart',
    value: 'InQuart',
  },
  {
    label: 'OutQuart',
    value: 'OutQuart',
  },
  {
    label: 'InOutQuart',
    value: 'InOutQuart',
  },
  {
    label: 'InQuint',
    value: 'InQuint',
  },
  {
    label: 'OutQuint',
    value: 'OutQuint',
  },
  {
    label: 'InOutQuint',
    value: 'InOutQuint',
  },
  {
    label: 'InSine',
    value: 'InSine',
  },
  {
    label: 'OutSine',
    value: 'OutSine',
  },
  {
    label: 'InOutSine',
    value: 'InOutSine',
  },
  {
    label: 'InExpo',
    value: 'InExpo',
  },
  {
    label: 'OutExpo',
    value: 'OutExpo',
  },
  {
    label: 'InOutExpo',
    value: 'InOutExpo',
  },
  {
    label: 'InCirc',
    value: 'InCirc',
  },
  {
    label: 'OutCirc',
    value: 'OutCirc',
  },
  {
    label: 'InOutCirc',
    value: 'InOutCirc',
  },
  {
    label: 'InBack',
    value: 'InBack',
  },
  {
    label: 'OutBack',
    value: 'OutBack',
  },
  {
    label: 'InOutBack',
    value: 'InOutBack',
  },
  {
    label: 'InElastic',
    value: 'InElastic',
  },
  {
    label: 'OutElastic',
    value: 'OutElastic',
  },
  {
    label: 'InOutElastic',
    value: 'InOutElastic',
  },
];

export const easeFunctionCategories: Array<SelectableValue<string>> = [
  {
    label: 'InOut easing',
    value: 'InOut*',
  },
  {
    label: 'Easing In',
    value: 'In[!O]*',
  },
  {
    label: 'Easing Out',
    value: 'Out*',
  },

  {
    label: 'Quadratic Functions',
    value: '*Quad',
  },
  {
    label: 'Cubic Functions',
    value: '*Cubic',
  },
  {
    label: 'Quart Functions',
    value: '*Quart',
  },
  {
    label: 'Quint Functions',
    value: '*Quint',
  },
  {
    label: 'Sine Functions',
    value: '*Sine',
  },
  {
    label: 'Exponential Functions',
    value: '*Expo',
  },
  {
    label: 'Circ Functions',
    value: '*Circ',
  },
  {
    label: 'Backoff Functions',
    value: '*Back',
  },
  {
    label: 'Elastic Functions',
    value: '*Elastic',
  },
];
