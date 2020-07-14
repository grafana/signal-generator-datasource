import { DataQuery, DataSourceJsonData } from '@grafana/data';

export enum DroneQueryType {
  Users = 'users',
  Repos = 'repos',
  Builds = 'builds',
  Logs = 'logs',
  Incomplete = 'incomplete',
  Nodes = 'Nodes',
  Servers = 'Servers',
}

export interface DroneQuery extends DataQuery {
  namespace?: string;
  name?: string;
  owner?: string;
  branch?: string;

  build?: number;
  stage?: number;
  step?: number;

  page?: number;
  pageSize?: number;
}

export interface WaveformDatasourceOptions extends DataSourceJsonData {
  // URL is used as host
}

export interface DroneSecureJsonData {
  token?: string;
}
