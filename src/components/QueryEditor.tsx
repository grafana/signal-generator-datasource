import React, { PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { Select, InlineField, InlineSwitch, InlineFieldRow } from '@grafana/ui';
import { DataSource } from '../DataSource';
import { SignalQuery, SignalDatasourceOptions, QueryType, SignalField, TimeFieldConfig } from '../types';
import { defaultSignal } from '../info';
import { SignalFieldEditor } from './SignalFieldEditor';
import { TimeFieldEditor } from './TimeFieldEditor';

type Props = QueryEditorProps<DataSource, SignalQuery, SignalDatasourceOptions>;

const queryTypes = [
  { label: 'Signal Generator', value: QueryType.AWG },
  { label: 'Variables', value: QueryType.Easing },
] as Array<SelectableValue<QueryType>>;

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

  onTimeChange = (time: TimeFieldConfig) => {
    const { onChange, query, onRunQuery } = this.props;
    const signal = { ...query.signal!, time };
    onChange({ ...query, signal });
    onRunQuery();
  };

  onSignalFieldChange = (v: SignalField | undefined, index: number, skipQuery?: boolean) => {
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
    if (!skipQuery) {
      onRunQuery();
    }
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

  onToggleStream = () => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, stream: !query.stream });
    onRunQuery();
  };

  renderQuery(query: SignalQuery) {
    if (query.queryType !== QueryType.AWG) {
      return <div>TODO: not implemented yet: ${query.queryType}</div>;
    }

    let signal = query.signal || defaultSignal;
    if (!signal.fields.length) {
      signal.fields = [...defaultSignal.fields];
    }
    return (
      <>
        <TimeFieldEditor time={signal.time} onChange={this.onTimeChange} />
        {signal.fields.map((s, idx) => {
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
        })}

        <InlineFieldRow>
          <InlineField label="Stream" labelWidth={8}>
            <InlineSwitch css="" onChange={this.onToggleStream} value={query.stream} />
          </InlineField>
        </InlineFieldRow>
      </>
    );
  }

  // renderEasing() {
  //   const { query } = this.props;
  //   const options = [...easeFunctionCategories, ...easeFunctions];
  //   const current = options.find(f => f.value === query.ease);

  //   return (
  //     <div className="gf-form">
  //       <InlineField label="Function" labelWidth={10} grow={true}>
  //         <Select
  //           options={options}
  //           value={current}
  //           onChange={this.onEaseChange}
  //           allowCustomValue={true}
  //           isClearable={true}
  //           isSearchable={true}
  //           placeholder="Show all functions"
  //           menuPlacement="bottom"
  //         />
  //       </InlineField>
  //     </div>
  //   );
  // }

  render() {
    const { query } = this.props;

    return (
      <>
        <div className="gf-form">
          <InlineField label="Query" labelWidth={8} grow={true}>
            <Select
              options={queryTypes}
              value={queryTypes.find((v) => v.value === query.queryType)}
              onChange={this.onQueryTypeChange}
              placeholder="Select query type"
              menuPlacement="bottom"
            />
          </InlineField>
        </div>
        {this.renderQuery(query)}
      </>
    );
  }
}
