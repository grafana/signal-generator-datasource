import React, { PureComponent } from 'react';
import { LegacyForms, InlineFormLabel } from '@grafana/ui';
const { Input, SecretFormField } = LegacyForms;
import {
  DataSourcePluginOptionsEditorProps,
  onUpdateDatasourceSecureJsonDataOption,
  updateDatasourcePluginResetOption,
} from '@grafana/data';
import { WaveformDatasourceOptions, DroneSecureJsonData } from '../types';

export type Props = DataSourcePluginOptionsEditorProps<WaveformDatasourceOptions, DroneSecureJsonData>;

export class ConfigEditor extends PureComponent<Props> {
  constructor(props: Props) {
    super(props);
  }

  onUpdateURL = (e: React.SyntheticEvent<HTMLInputElement>) => {
    const { options, onOptionsChange } = this.props;
    onOptionsChange({
      ...options,
      url: e.currentTarget.value,
      access: 'proxy',
    });
  };

  render() {
    const { options } = this.props;
    const { secureJsonFields } = options;
    const secureJsonData = (options.secureJsonData || {}) as DroneSecureJsonData;

    return (
      <>
        <h3 className="page-heading">Connection</h3>
        <div className="gf-form-inline">
          <div className="gf-form">
            <InlineFormLabel
              className="width-10"
              tooltip="This URL needs to be accessible from the grafana backend/server."
            >
              URL
            </InlineFormLabel>
            <div className="width-20">
              <Input
                className="width-20"
                value={options.url || ''}
                placeholder="http://drone.company.com"
                onChange={this.onUpdateURL}
              />
            </div>
          </div>
        </div>
        <div className="gf-form-inline">
          <div className="gf-form">
            <SecretFormField
              isConfigured={(secureJsonFields && secureJsonFields.token) as boolean}
              value={secureJsonData.token || ''}
              label="Token"
              placeholder="service accout token"
              labelWidth={10}
              inputWidth={20}
              onReset={() => updateDatasourcePluginResetOption(this.props, 'token')}
              onChange={onUpdateDatasourceSecureJsonDataOption(this.props, 'token')}
            />
          </div>
        </div>
      </>
    );
  }
}
