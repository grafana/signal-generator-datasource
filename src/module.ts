import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import { ConfigEditor, QueryEditor } from './components';
import { AWGQuery, AWGDatasourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, AWGQuery, AWGDatasourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
