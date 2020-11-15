import { SelectableValue } from '@grafana/data';
import { SignalConfig } from 'types';

export const defaultSignal: SignalConfig = {
  time: {
    period: '1m',
  },
  fields: [{ name: 'Value', expr: 'Sine(x)' }],
};

export const standardWaves: Array<SelectableValue<string>> = [
  {
    label: 'Sine(x)',
    value: 'Sine(x)',
    description: 'A curve that describes a smooth periodic oscillation',
  },
  {
    label: 'Triangle(x)',
    value: 'Triangle(x)',
    description: 'triangle qave',
  },
  {
    label: 'Square(x)',
    value: 'Square(x)',
    description: 'square wave with a 50% duty cycle',
  },
  {
    label: 'Square(x,0.2)',
    value: 'Square(x,0.2)',
    description: 'square wave with a 20% duty cycle',
  },
  {
    label: 'Square(x,0.8)',
    value: 'Square(x,0.8)',
    description: 'square wave with an 80% duty cycle',
  },
  {
    label: 'Sawtooth(x)',
    value: 'Sawtooth(x)',
    description: 'square wave with a 50% duty cycle',
  },
  {
    label: 'Noise()',
    value: 'Noise()',
    description: 'random values between +/- 1',
  },
  {
    label: 'CSV(x,1,3,-1,-3)',
    value: 'CSV(x,1,3,-1,-3)',
    description: 'csv wave',
  },
  {
    label: 'x',
    value: 'x',
    description: 'value that ramps between -Pi and +Pi depending on the period',
  },
];

// export const waveformTypes: Array<SelectableValue<WaveformType>> = [
//   {
//     label: 'Sine Wave',
//     value: WaveformType.Sin,
//     description: 'A curve that describes a smooth periodic oscillation',
//   },
//   {
//     label: 'Square Wave',
//     value: WaveformType.Square,
//     description: 'fixed between minimum and maximum values',
//   },
//   {
//     label: 'Triangle Wave',
//     value: WaveformType.Triangle,
//     description: 'fixed between minimum and maximum values',
//   },
//   {
//     label: 'Sawtooth Wave',
//     value: WaveformType.Sawtooth,
//     description: 'fixed between minimum and maximum values',
//   },
//   {
//     label: 'Random noise',
//     value: WaveformType.Noise,
//     description: 'random values',
//   },
//   {
//     label: 'CSV Values',
//     value: WaveformType.CSV,
//     description: 'Animated values',
//   },
//   {
//     label: 'Calculation',
//     value: WaveformType.Calculation,
//     description: 'Calculate a value',
//   },
// ];

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
