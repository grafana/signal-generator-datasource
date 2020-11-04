import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';

import { SignalQuery, SignalDatasourceOptions } from './types';

export class DataSource extends DataSourceWithBackend<SignalQuery, SignalDatasourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<SignalDatasourceOptions>) {
    super(instanceSettings);
  }
}
