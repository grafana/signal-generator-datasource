import { DataQuery, DataSourceJsonData } from '@grafana/data';

export enum AWGQueryType {
  Signal = 'signal',
  Stream = 'stream',
}

export interface AWGQuery extends DataQuery {
  queryType?: AWGQueryType;
}

export interface AWGDatasourceOptions extends DataSourceJsonData {
  // URL is used as host
}

export interface AWGSecureJsonData {
  // nothing for now
}
