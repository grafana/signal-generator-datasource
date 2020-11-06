import React, { PureComponent } from 'react';
import { SelectableValue } from '@grafana/data';
import { Input, Select, InlineField } from '@grafana/ui';
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
    const copy = { ...wave, type: sel.value! };
    if (copy.type === WaveformType.CSV) {
      if (!copy.args) {
        copy.args = 'InOutQuad';
      }
      if (!copy.points) {
        copy.points = [1, 0.5, 7, 3];
      }
    } else if (copy.type === WaveformType.Calculation) {
      copy.args = 'sin(time)/time';
    }

    if (copy.type === WaveformType.Square) {
      copy.duty = 0.5;
    } else {
      delete copy.duty;
    }
    onChange(copy, index);
  };

  onPeriodChange = (sel: SelectableValue<string>) => {
    const { onChange, wave, index } = this.props;
    onChange({ ...wave, period: sel.value! }, index);
  };

  onAmplitudeChange = (v: React.SyntheticEvent<HTMLInputElement>) => {
    const { onChange, wave, index } = this.props;
    const amplitude = v.currentTarget.valueAsNumber;
    onChange({ ...wave, amplitude }, index);
  };

  onOffsetChange = (v: React.SyntheticEvent<HTMLInputElement>) => {
    const { onChange, wave, index } = this.props;
    const offset = v.currentTarget.valueAsNumber;
    onChange({ ...wave, offset }, index);
  };

  onPhaseChange = (v: React.SyntheticEvent<HTMLInputElement>) => {
    const { onChange, wave, index } = this.props;
    const phase = v.currentTarget.valueAsNumber;
    onChange({ ...wave, phase }, index);
  };

  onArgsChange = (v: React.SyntheticEvent<HTMLInputElement>) => {
    const { onChange, wave, index } = this.props;
    const args = v.currentTarget.value;
    onChange({ ...wave, args }, index);
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

    if (wave.type === WaveformType.Calculation) {
      return (
        <>
          <div className="gf-form">
            <InlineField label=">>>" labelWidth={10}>
              <Select
                options={waveformTypes}
                value={waveformTypes.find(v => v.value === wave.type)}
                onChange={this.onQueryTypeChange}
                placeholder="Select waveform"
                menuPlacement="bottom"
              />
            </InlineField>
            <InlineField label="Calculation" grow={true}>
              <Input css="" value={wave.args || ''} onChange={this.onArgsChange} />
            </InlineField>
          </div>
        </>
      );
    }

    return (
      <>
        <div className="gf-form">
          <InlineField label=">>>" labelWidth={10} grow={true}>
            <Select
              options={waveformTypes}
              value={waveformTypes.find(v => v.value === wave.type)}
              onChange={this.onQueryTypeChange}
              placeholder="Select waveform"
              menuPlacement="bottom"
            />
          </InlineField>
          <InlineField label="Amplitude">
            <Input
              css=""
              type="number"
              width={6}
              step={0.1}
              defaultValue={wave.amplitude || 1}
              onBlur={this.onAmplitudeChange}
            />
          </InlineField>
          <InlineField label="Offset">
            <Input
              css=""
              type="number"
              width={6}
              step={0.1}
              defaultValue={wave.offset || 0}
              onBlur={this.onOffsetChange}
            />
          </InlineField>
          <InlineField label="Phase">
            <Input
              css=""
              type="number"
              width={6}
              step={0.1}
              min={0}
              max={1}
              defaultValue={wave.phase || 0}
              onBlur={this.onPhaseChange}
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
        </div>
      </>
    );
  }
}
