import { DataQuery, DataSourceJsonData } from '@grafana/data';

export enum QueryType {
  AWG = 'AWG',
  Easing = 'easing',
}

export enum WaveformType {
  Sin = 'Sin',
  Pulse = 'Pulse',
  Sawtooth = 'Sawtooth',
  Sinc = 'Sinc',
  Noise = 'Noise',
  CSV = 'CSV',
}

export interface WaveformArgs {
  type: WaveformType;
  period: number; // in seconds
  amplitude: number;
  duty?: number;
  points?: number[]; // for CSV
  ease?: string; // Ease function
}

export interface SignalQuery extends DataQuery {
  queryType?: QueryType;
}

export interface SignalDatasourceOptions extends DataSourceJsonData {
  // nothing for now
}

export interface SignalSecureJsonData {
  // nothing for now
}
