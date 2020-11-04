import React, { PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { Select, InlineField } from '@grafana/ui';
import { DataSource } from '../DataSource';
import { SignalQuery, SignalDatasourceOptions, QueryType } from '../types';
import { easeFunctionCategories, easeFunctions } from '../info';

type Props = QueryEditorProps<DataSource, SignalQuery, SignalDatasourceOptions>;

const queryTypes = [
  { label: 'Waveform', value: QueryType.AWG },
  { label: 'Easing', value: QueryType.Easing },
] as Array<SelectableValue<QueryType>>;

export class QueryEditor extends PureComponent<Props> {
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

  renderAWG() {
    return <div>TODO...</div>;
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

    if (!query.queryType) {
      query.queryType = QueryType.AWG;
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
            />
          </InlineField>
        </div>
        {query.queryType === QueryType.AWG && this.renderAWG()}
        {query.queryType === QueryType.Easing && this.renderEasing()}
      </>
    );
  }
}
