import React, { PureComponent } from 'react';
import { DataSourcePluginOptionsEditorProps, onUpdateDatasourceJsonDataOption } from '@grafana/data';
import { SignalDatasourceOptions, SignalSecureJsonData } from '../types';

import { Field, Input, TextArea } from '@grafana/ui';

export type Props = DataSourcePluginOptionsEditorProps<SignalDatasourceOptions, SignalSecureJsonData>;

export class ConfigEditor extends PureComponent<Props> {
  constructor(props: Props) {
    super(props);
  }

  render() {
    const { options } = this.props;

    return (
      <>
        <div>
          <Field label="Grafana URL (temporary for streaming)">
            <Input
              value={options.jsonData.live}
              placeholder="http://localhost:3000/"
              css=""
              autoComplete="off"
              onChange={onUpdateDatasourceJsonDataOption(this.props, 'live')}
            />
          </Field>

          <Field label="Capture">
            <TextArea
              value={options.jsonData.captureX}
              placeholder="/path/to/tags.json"
              css=""
              autoComplete="off"
              rows={4}
              onChange={onUpdateDatasourceJsonDataOption(this.props, 'captureX') as any}
            />
          </Field>
        </div>
      </>
    );
  }
}
