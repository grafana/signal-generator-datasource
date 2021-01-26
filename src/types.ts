import { DataQuery, DataSourceJsonData, FieldConfig, Labels } from '@grafana/data';

export enum QueryType {
  AWG = 'AWG',
  Easing = 'easing',
  Streams = 'streams',
}

// export enum WaveformType {
//   Sin = 'Sin',
//   Square = 'Square',
//   Triangle = 'Triangle',
//   Sawtooth = 'Sawtooth',
//   Noise = 'Noise',
//   CSV = 'CSV',
//   Calculation = 'Calculation',
// }

export interface SignalField {
  name?: string;
  expr: string;
  config?: FieldConfig;
  labels?: Labels;
}

export interface TimeFieldConfig {
  period: string;
}
export interface RangeFieldConfig {
  min: number;
  max: number;
  ease?: string; // ease function matcher
}

export interface SignalConfig {
  name?: string;
  time: TimeFieldConfig;
  fields: SignalField[];
}

/**
 * Metadata attached to DataFrame results
 */
export interface SignalCustomMeta {
  streamKey?: string;
}

export interface SignalQuery extends DataQuery {
  queryType?: QueryType;
  signal?: SignalConfig;
}

export interface SignalDatasourceOptions extends DataSourceJsonData {
  live: string;
  captureX: string; // paths to local files
}

export interface SignalSecureJsonData {
  // nothing for now
}
