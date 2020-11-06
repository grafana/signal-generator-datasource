import React, { PureComponent } from 'react';
import { InlineField, Input, UnitPicker } from '@grafana/ui';
import { SignalArgs, WaveformArgs } from '../types';
import { defaultWave } from '../info';
import { WaveEditor } from './WaveEditor';
import { formatLabels, parseLabels } from '@grafana/data';

interface Props {
  signal: SignalArgs;
  index: number;
  onChange: (value: SignalArgs, index: number) => void;
}

export class SignalEditor extends PureComponent<Props> {
  onWaveChange = (wave: WaveformArgs | undefined, index: number) => {
    const { onChange, signal } = this.props;
    const copy = [...signal.component];
    if (wave) {
      copy[index] = wave;
    } else {
      // Remove the value
      copy.splice(index, 1);
    }
    onChange(
      {
        ...signal,
        component: copy,
      },
      this.props.index
    );
  };

  onNameChange = (v: React.SyntheticEvent<HTMLInputElement>) => {
    const { onChange, signal, index } = this.props;
    const name = v.currentTarget.value;
    onChange({ ...signal, name }, index);
  };

  onLabelsChanged = (v: React.SyntheticEvent<HTMLInputElement>) => {
    const { onChange, signal, index } = this.props;
    const txt = v.currentTarget.value;
    const labels = txt ? parseLabels(txt) : undefined;
    onChange({ ...signal, labels }, index);
  };

  onUnitChanged = (v?: string) => {
    const { onChange, signal, index } = this.props;
    const copy = { ...signal };
    if (v) {
      if (!copy.config) {
        copy.config = {};
      }
      copy.config.unit = v;
    } else if (copy.config) {
      delete copy.config.unit;
    }
    onChange(copy, index);
  };

  render() {
    const { signal } = this.props;
    if (!signal.component) {
      signal.component = [{ ...defaultWave }];
    }

    return (
      <>
        <div className="gf-form">
          <InlineField label="Signal" labelWidth={10} grow={true}>
            <Input css="" value={signal.name || ''} onChange={this.onNameChange} placeholder="Field name" />
          </InlineField>
          <InlineField label="Labels">
            <Input
              css=""
              defaultValue={formatLabels(signal.labels!)}
              onBlur={this.onLabelsChanged}
              placeholder="labels"
            />
          </InlineField>
          <InlineField label="Unit">
            <UnitPicker value={signal.config?.unit} onChange={this.onUnitChanged} width={15} />
          </InlineField>
        </div>
        {signal.component.map((w, idx) => {
          return <WaveEditor wave={w} index={idx} key={idx} onChange={this.onWaveChange} />;
        })}
      </>
    );
  }
}
