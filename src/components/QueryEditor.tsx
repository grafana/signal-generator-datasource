import React, { PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { Select, InlineFormLabel } from '@grafana/ui';
import { DataSource } from '../DataSource';
import { AWGQuery, AWGDatasourceOptions, AWGQueryType } from '../types';

type Props = QueryEditorProps<DataSource, AWGQuery, AWGDatasourceOptions>;

const queryTypes = [
  { label: 'Random Signal', value: AWGQueryType.Signal },
  { label: 'Stream', value: AWGQueryType.Stream },
] as Array<SelectableValue<AWGQueryType>>;

// const defaultQuery: KeyValue<Partial<AWGQuery>> = {
//   [AWGQueryType.Signal]: {},
//   [AWGQueryType.Stream]: {},
// };

export class QueryEditor extends PureComponent<Props> {
  onQueryTypeChange = (val: SelectableValue) => {
    const { query, onChange } = this.props;
    onChange({
      ...query,
      queryType: val.value,
    });
    this.props.onRunQuery();
  };

  // onNamespaceChange = (event: ChangeEvent<HTMLInputElement>) => {
  //   this.props.onChange({
  //     ...this.props.query,
  //     namespace: event.target.value,
  //   });
  //   this.props.onRunQuery();
  // };

  // onNameChange = (event: ChangeEvent<HTMLInputElement>) => {
  //   this.props.onChange({
  //     ...this.props.query,
  //     name: event.target.value,
  //   });
  //   this.props.onRunQuery();
  // };

  // onBuildChange = (event: ChangeEvent<HTMLInputElement>) => {
  //   this.props.onChange({
  //     ...this.props.query,
  //     build: +event.target.value,
  //   });
  //   this.props.onRunQuery();
  // };
  // onStageChange = (event: ChangeEvent<HTMLInputElement>) => {
  //   this.props.onChange({
  //     ...this.props.query,
  //     stage: +event.target.value,
  //   });
  //   this.props.onRunQuery();
  // };
  // onStepChange = (event: ChangeEvent<HTMLInputElement>) => {
  //   this.props.onChange({
  //     ...this.props.query,
  //     step: +event.target.value,
  //   });
  //   this.props.onRunQuery();
  // };

  render() {
    const { query } = this.props;
    const labelWidth = 8;

    if (!query.queryType) {
      query.queryType = AWGQueryType.Signal;
    }
    //  const standard = defaultQuery[query.queryType] || {};

    return (
      <>
        <div className="gf-form-inline">
          <div className="gf-form">
            <InlineFormLabel width={labelWidth} className="query-keyword">
              Query
            </InlineFormLabel>
            <Select
              className="width-20"
              value={queryTypes.find(queryType => queryType.value === query.queryType)}
              options={queryTypes}
              onChange={this.onQueryTypeChange}
            />
          </div>
          <div className="gf-form--grow">
            <label className="gf-form-label gf-form-label--grow"></label>
          </div>
        </div>
      </>
    );
  }
}
