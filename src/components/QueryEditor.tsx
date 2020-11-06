import React, { PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { Select, InlineField } from '@grafana/ui';
import { DataSource } from '../DataSource';
import { SignalQuery, SignalDatasourceOptions, QueryType, SignalArgs } from '../types';
import { defaultSignal, easeFunctionCategories, easeFunctions } from '../info';
import { SignalEditor } from './SignalEditor';

type Props = QueryEditorProps<DataSource, SignalQuery, SignalDatasourceOptions>;

const queryTypes = [
  { label: 'Waveform', value: QueryType.AWG },
  { label: 'Easing', value: QueryType.Easing },
] as Array<SelectableValue<QueryType>>;

export class QueryEditor extends PureComponent<Props> {
  componentDidMount() {
    const { onChange, query } = this.props;

    let changed = false;
    if (!query.queryType) {
      query.queryType = QueryType.AWG;
      changed = true;
    }
    if (!query.signals) {
      query.signals = [{ ...defaultSignal }];
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

  onSignalChange = (v: SignalArgs | undefined, index: number) => {
    const { onChange, query, onRunQuery } = this.props;
    const copy = [...query.signals];
    if (v) {
      copy[index] = v;
    } else {
      // Remove the value
      copy.splice(index, 1);
    }
    onChange({ ...query, signals: copy });
    onRunQuery();
  };

  renderAWG() {
    const { query } = this.props;
    let signals = query.signals;
    if (!signals || !signals.length) {
      signals = [{ ...defaultSignal }];
    }
    return signals.map((s, idx) => {
      return <SignalEditor signal={s} index={idx} key={idx} onChange={this.onSignalChange} />;
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

    return (
      <>
        <div className="gf-form">
          <InlineField label="Query type" labelWidth={10} grow={true}>
            <Select
              options={queryTypes}
              value={queryTypes.find(v => v.value === query.queryType)}
              onChange={this.onQueryTypeChange}
              placeholder="Select query type"
            />
          </InlineField>
        </div>
        {query.queryType === QueryType.AWG && this.renderAWG()}
        {query.queryType === QueryType.Easing && this.renderEasing()}
      </>
    );
  }
}
