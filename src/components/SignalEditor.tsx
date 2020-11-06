import React, { PureComponent } from 'react';
import { InlineField } from '@grafana/ui';
import { SignalArgs, WaveformArgs } from '../types';
import { defaultWave } from '../info';
import { WaveEditor } from './WaveEditor';

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

  render() {
    const { signal } = this.props;
    if (!signal.component) {
      signal.component = [{ ...defaultWave }];
    }

    return (
      <>
        <div className="gf-form">
          <InlineField label="Signal" labelWidth={10} grow={true}>
            <div>NAME</div>
          </InlineField>
          <InlineField label="Labels">
            <div>LABELS</div>
          </InlineField>
          <InlineField label="Unit">
            <div>UNIT</div>
          </InlineField>
        </div>
        {signal.component.map((w, idx) => {
          return <WaveEditor wave={w} index={idx} key={idx} onChange={this.onWaveChange} />;
        })}
      </>
    );
  }
}
