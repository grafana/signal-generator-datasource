import React, { PureComponent } from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { SignalDatasourceOptions, SignalSecureJsonData } from '../types';

export type Props = DataSourcePluginOptionsEditorProps<SignalDatasourceOptions, SignalSecureJsonData>;

export class ConfigEditor extends PureComponent<Props> {
  constructor(props: Props) {
    super(props);
  }

  render() {
    return (
      <>
        <div>Signal generator</div>
      </>
    );
  }
}
