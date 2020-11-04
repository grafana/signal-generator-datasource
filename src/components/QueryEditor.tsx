import React, { PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { Select, InlineField } from '@grafana/ui';
import { DataSource } from '../DataSource';
import { SignalQuery, SignalDatasourceOptions, QueryType, WaveformType, WaveformArgs } from '../types';
import { easeFunctionCategories, easeFunctions } from '../info';
import { WaveEditor } from './WaveEditor';

type Props = QueryEditorProps<DataSource, SignalQuery, SignalDatasourceOptions>;

const queryTypes = [
  { label: 'Waveform', value: QueryType.AWG },
  { label: 'Easing', value: QueryType.Easing },
] as Array<SelectableValue<QueryType>>;

const defaultWave: WaveformArgs = {
  type: WaveformType.Sin,
  period: '1m',
  amplitude: 1,
};

export class QueryEditor extends PureComponent<Props> {
  componentDidMount() {
    const { onChange, query } = this.props;

    let changed = false;
    if (!query.queryType) {
      query.queryType = QueryType.AWG;
      changed = true;
    }
    if (!query.wave) {
      query.wave = [{ ...defaultWave }];
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

  onWaveChange = (wave: WaveformArgs | undefined, index: number) => {
    const { onChange, query, onRunQuery } = this.props;
    const copy = [...query.wave];
    if (wave) {
      copy[index] = wave;
    } else {
      // Remove the value
      copy.splice(index, 1);
    }
    onChange({ ...query, wave: copy });
    onRunQuery();
  };

  renderAWG() {
    const { query } = this.props;
    let waves = query.wave;
    if (!waves || !waves.length) {
      waves = [{ ...defaultWave }];
    }
    return waves.map((wave, idx) => {
      return <WaveEditor wave={wave} index={idx} key={idx} onChange={this.onWaveChange} />;
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
