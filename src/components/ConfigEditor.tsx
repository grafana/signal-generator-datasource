import React, { PureComponent } from 'react';
import { LegacyForms, InlineFormLabel } from '@grafana/ui';
const { Input } = LegacyForms;
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { AWGDatasourceOptions, AWGSecureJsonData } from '../types';

export type Props = DataSourcePluginOptionsEditorProps<AWGDatasourceOptions, AWGSecureJsonData>;

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
    if (true) {
      return <div />;
    }

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
      </>
    );
  }
}
