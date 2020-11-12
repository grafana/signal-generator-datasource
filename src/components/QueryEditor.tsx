import React, { PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { Select, InlineField } from '@grafana/ui';
import { DataSource } from '../DataSource';
import { SignalQuery, SignalDatasourceOptions, QueryType, SignalField } from '../types';
import { defaultSignal, easeFunctionCategories, easeFunctions } from '../info';
import { SignalFieldEditor } from './SignalFieldEditor';

type Props = QueryEditorProps<DataSource, SignalQuery, SignalDatasourceOptions>;

const queryTypes = [
  { label: 'Waveform', value: QueryType.AWG },
  { label: 'Easing', value: QueryType.Easing },
] as Array<SelectableValue<QueryType>>;

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

export class QueryEditor extends PureComponent<Props> {
  componentDidMount() {
    const { onChange, query } = this.props;

    let changed = false;
    if (!query.queryType) {
      query.queryType = QueryType.AWG;
      changed = true;
    }
    if (!query.signal) {
      query.signal = { ...defaultSignal };
      changed = true;
    }
    if (changed) {
      onChange({ ...query });
    }
  }

  onQueryTypeChange = (sel: SelectableValue<QueryType>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, queryType: sel.value });
    onRunQuery();
  };

  onEaseChange = (sel: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, ease: sel.value });
    onRunQuery();
  };

  onPeriodChange = (sel: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, period: sel.value! });
    onRunQuery();
  };

  onSignalFieldChange = (v: SignalField | undefined, index: number) => {
    const { onChange, query, onRunQuery } = this.props;
    const signal = { ...query.signal! };
    const fields = [...signal.fields];
    if (v) {
      fields[index] = v;
    } else {
      // Remove the value
      fields.splice(index, 1);
    }
    signal.fields = fields;

    onChange({ ...query, signal });
    onRunQuery();
  };

  onAddExpr = () => {
    const { onChange, query, onRunQuery } = this.props;
    let { signal } = query;
    if (!signal) {
      signal = { ...defaultSignal };
    } else {
      const fields = [...signal.fields, { ...defaultSignal.fields[0] }];
      signal = { ...signal, fields };
    }
    onChange({ ...query, signal });
    onRunQuery();
  };

  renderAWG() {
    const { query } = this.props;
    let signal = query.signal || defaultSignal;
    if (!signal.fields.length) {
      signal.fields = [...defaultSignal.fields];
    }
    return signal.fields.map((s, idx) => {
      const isLast = idx === signal.fields.length - 1;
      return (
        <SignalFieldEditor
          signal={s}
          index={idx}
          key={idx}
          onChange={this.onSignalFieldChange}
          onAddExpr={isLast ? this.onAddExpr : undefined}
        />
      );
    });
  }

  renderEasing() {
    const { query } = this.props;
    const options = [...easeFunctionCategories, ...easeFunctions];
    const current = options.find(f => f.value === query.ease);

    return (
      <div className="gf-form">
        <InlineField label="Function" labelWidth={10} grow={true}>
          <Select
            options={options}
            value={current}
            onChange={this.onEaseChange}
            allowCustomValue={true}
            isClearable={true}
            isSearchable={true}
            placeholder="Show all functions"
            menuPlacement="bottom"
          />
        </InlineField>
      </div>
    );
  }

  render() {
    const { query } = this.props;

    const periods = [...commonPeriods];
    let period = periods.find(p => p.value === query?.period);
    if (!period && query?.period) {
      period = {
        label: query.period,
        value: query.period,
      };
      periods.push(period);
    }

    return (
      <>
        <div className="gf-form">
          <InlineField label="Query type" labelWidth={10} grow={true}>
            <Select
              options={queryTypes}
              value={queryTypes.find(v => v.value === query.queryType)}
              onChange={this.onQueryTypeChange}
              placeholder="Select query type"
              menuPlacement="bottom"
            />
          </InlineField>
          {query.queryType === QueryType.AWG && (
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
          )}
        </div>
        {query.queryType === QueryType.AWG && this.renderAWG()}
        {query.queryType === QueryType.Easing && this.renderEasing()}
      </>
    );
  }
}
