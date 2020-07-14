import React, { PureComponent, ChangeEvent } from 'react';
import { QueryEditorProps, SelectableValue, KeyValue } from '@grafana/data';
import { Select, InlineFormLabel } from '@grafana/ui';
import { DataSource } from '../DataSource';
import { DroneQuery, WaveformDatasourceOptions, DroneQueryType } from '../types';

type Props = QueryEditorProps<DataSource, DroneQuery, WaveformDatasourceOptions>;

const queryTypes = [
  { label: 'Show Repositories', value: DroneQueryType.Repos },
  { label: 'Show Builds', value: DroneQueryType.Builds },
  { label: 'Show Logs', value: DroneQueryType.Logs },
] as Array<SelectableValue<DroneQueryType>>;

const defaultQuery: KeyValue<Partial<DroneQuery>> = {
  [DroneQueryType.Repos]: {},
  [DroneQueryType.Builds]: {
    namespace: 'grafana',
    name: 'grafana',
    pageSize: 10,
  },
  [DroneQueryType.Logs]: {
    namespace: 'grafana',
    name: 'grafana',
    build: 23,
    stage: 1,
    step: 1,
  },
};

export class QueryEditor extends PureComponent<Props> {
  onQueryTypeChange = (val: SelectableValue) => {
    const { query, onChange } = this.props;
    onChange({
      ...query,
      queryType: val.value,
    });
    this.props.onRunQuery();
  };

  onNamespaceChange = (event: ChangeEvent<HTMLInputElement>) => {
    this.props.onChange({
      ...this.props.query,
      namespace: event.target.value,
    });
    this.props.onRunQuery();
  };

  onNameChange = (event: ChangeEvent<HTMLInputElement>) => {
    this.props.onChange({
      ...this.props.query,
      name: event.target.value,
    });
    this.props.onRunQuery();
  };

  onBuildChange = (event: ChangeEvent<HTMLInputElement>) => {
    this.props.onChange({
      ...this.props.query,
      build: +event.target.value,
    });
    this.props.onRunQuery();
  };
  onStageChange = (event: ChangeEvent<HTMLInputElement>) => {
    this.props.onChange({
      ...this.props.query,
      stage: +event.target.value,
    });
    this.props.onRunQuery();
  };
  onStepChange = (event: ChangeEvent<HTMLInputElement>) => {
    this.props.onChange({
      ...this.props.query,
      step: +event.target.value,
    });
    this.props.onRunQuery();
  };

  render() {
    const { query, onRunQuery } = this.props;
    const labelWidth = 8;

    if (!query.queryType) {
      query.queryType = DroneQueryType.Repos;
    }
    const standard = defaultQuery[query.queryType] || {};

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

        {standard.namespace && (
          <div className="gf-form-inline">
            <InlineFormLabel width={labelWidth} className="query-keyword">
              Namespace
            </InlineFormLabel>
            <input
              className="gf-form-input width-14"
              value={query.namespace || ''}
              placeholder={standard.namespace}
              onChange={this.onNamespaceChange}
              onBlur={onRunQuery}
            ></input>
            <div className="gf-form gf-form--grow">
              <div className="gf-form-label gf-form-label--grow" />
            </div>
          </div>
        )}

        {standard.name && (
          <div className="gf-form-inline">
            <InlineFormLabel width={labelWidth} className="query-keyword">
              Name
            </InlineFormLabel>
            <input
              className="gf-form-input width-14"
              value={query.name || ''}
              placeholder={standard.name}
              onChange={this.onNameChange}
              onBlur={onRunQuery}
            ></input>
            <div className="gf-form gf-form--grow">
              <div className="gf-form-label gf-form-label--grow" />
            </div>
          </div>
        )}

        {standard.build && (
          <div className="gf-form-inline">
            <InlineFormLabel width={labelWidth} className="query-keyword">
              Build Number
            </InlineFormLabel>
            <input
              type="number"
              className="gf-form-input width-14"
              value={query.build || ''}
              placeholder={standard.build + ''}
              onChange={this.onBuildChange}
              onBlur={onRunQuery}
            ></input>
            <div className="gf-form gf-form--grow">
              <div className="gf-form-label gf-form-label--grow" />
            </div>
          </div>
        )}

        {standard.stage && (
          <div className="gf-form-inline">
            <InlineFormLabel width={labelWidth} className="query-keyword">
              Stage Number
            </InlineFormLabel>
            <input
              type="number"
              className="gf-form-input width-14"
              value={query.stage || ''}
              placeholder={standard.stage + ''}
              onChange={this.onStageChange}
              onBlur={onRunQuery}
            ></input>
            <div className="gf-form gf-form--grow">
              <div className="gf-form-label gf-form-label--grow" />
            </div>
          </div>
        )}

        {standard.step && (
          <div className="gf-form-inline">
            <InlineFormLabel width={labelWidth} className="query-keyword">
              Step Number
            </InlineFormLabel>
            <input
              type="number"
              className="gf-form-input width-14"
              value={query.step || ''}
              placeholder={standard.step + ''}
              onChange={this.onStepChange}
              onBlur={onRunQuery}
            ></input>
            <div className="gf-form gf-form--grow">
              <div className="gf-form-label gf-form-label--grow" />
            </div>
          </div>
        )}
      </>
    );
  }
}
