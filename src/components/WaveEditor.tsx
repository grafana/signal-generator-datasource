import React, { PureComponent } from 'react';
import { SelectableValue } from '@grafana/data';
import { Select, InlineField } from '@grafana/ui';
import { WaveformArgs, WaveformType } from '../types';
import { waveformTypes } from '../info';

interface Props {
  wave: WaveformArgs;
  index: number;
  onChange: (value: WaveformArgs, index: number) => void;
}

export const commonPeriods: Array<SelectableValue<string>> = [
  {
    label: '1m',
    value: '1m',
  },
  {
    label: '10s',
    value: '10s',
  },
  {
    label: '1h',
    value: '1h',
  },
  {
    label: 'range/2',
    value: 'range/2',
  },
];

export class WaveEditor extends PureComponent<Props> {
  onQueryTypeChange = (sel: SelectableValue<WaveformType>) => {
    const { onChange, wave, index } = this.props;
    onChange({ ...wave, type: sel.value! }, index);
  };

  onPeriodChange = (sel: SelectableValue<string>) => {
    const { onChange, wave, index } = this.props;
    onChange({ ...wave, period: sel.value! }, index);
  };

  onAmplitudeChange = (v: any) => {
    // const { onChange, wave, index } = this.props;
    // onChange({ ...wave, amplitude: v }, index);
    console.log('AMP', v);
  };
  render() {
    const { wave } = this.props;
    const periods = [...commonPeriods];
    let period = periods.find(p => p.value === wave.period);
    if (!period && wave.period) {
      period = {
        label: wave.period,
        value: wave.period,
      };
      periods.push(period);
    }

    return (
      <>
        <div className="gf-form">
          <InlineField label="Query type" labelWidth={10} grow={true}>
            <Select
              options={waveformTypes}
              value={waveformTypes.find(v => v.value === wave.type)}
              onChange={this.onQueryTypeChange}
              placeholder="Select waveform"
              menuPlacement="bottom"
            />
          </InlineField>
          <InlineField label="Period">
            <Select
              options={periods}
              value={period}
              onChange={this.onPeriodChange}
              placeholder="Enter period"
              allowCustomValue={true}
              menuPlacement="bottom"
            />
          </InlineField>
          <InlineField label="Amplitude">
            <div>TODO</div>
          </InlineField>
        </div>
      </>
    );
  }
}
