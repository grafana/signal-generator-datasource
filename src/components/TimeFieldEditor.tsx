import React, { PureComponent } from 'react';
import { InlineField, InlineLabel, Select } from '@grafana/ui';
import { TimeFieldConfig } from '../types';
import { SelectableValue } from '@grafana/data';

interface Props {
  time: TimeFieldConfig;
  onChange: (time: TimeFieldConfig) => void;
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

export class TimeFieldEditor extends PureComponent<Props> {
  onPeriodChange = (sel: SelectableValue<string>) => {
    const { onChange, time } = this.props;
    onChange({ ...time, period: sel.value! });
  };

  // onNameChange = (v: React.SyntheticEvent<HTMLInputElement>) => {
  //   const { onChange, time } = this.props;
  //   const name = v.currentTarget.value;
  //   onChange({ ...time, name });
  // };

  // onLabelsChanged = (v: React.SyntheticEvent<HTMLInputElement>) => {
  //   const { onChange, time, index } = this.props;
  //   const txt = v.currentTarget.value;
  //   const labels = txt ? parseLabels(txt) : undefined;
  //   onChange({ ...time, labels }, index);
  // };

  render() {
    const { time } = this.props;
    const periods = [...commonPeriods];
    let period = periods.find((p) => p.value === time?.period);
    if (!period && time?.period) {
      period = {
        label: time.period,
        value: time.period,
      };
      periods.push(period);
    }

    return (
      <>
        <div className="gf-form">
          <InlineLabel width={8}>Time</InlineLabel>
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
          <InlineField label="Variables" tooltip="avaliable in the equations" grow={true}>
            <InlineLabel>x = 0-2pi, p = 0-1 percent</InlineLabel>
          </InlineField>
        </div>
      </>
    );
  }
}
