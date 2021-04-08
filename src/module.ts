import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/QueryEditor';
import { SignalQuery, SignalDatasourceOptions } from './types';
import { LiveMeasurementsSupport } from 'support';

export const plugin = new DataSourcePlugin<DataSource, SignalQuery, SignalDatasourceOptions>(DataSource)
  .setChannelSupport(new LiveMeasurementsSupport())
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
