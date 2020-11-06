import { SelectableValue } from '@grafana/data';
import { SignalArgs, WaveformArgs, WaveformType } from 'types';

export const defaultWave: WaveformArgs = {
  type: WaveformType.Sin,
  period: '1m',
  amplitude: 1,
  offset: 0,
  phase: 0,
};

export const defaultSignal: SignalArgs = {
  name: '',
  component: [{ ...defaultWave }],
};

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
    label: 'Random noise',
    value: WaveformType.Noise,
    description: 'random values',
  },
  {
    label: 'CSV Values',
    value: WaveformType.CSV,
    description: 'Animated values',
  },
  {
    label: 'Calculation',
    value: WaveformType.Calculation,
    description: 'Calculate a value',
  },
];

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
