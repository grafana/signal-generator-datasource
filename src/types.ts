import { DataQuery, DataSourceJsonData, FieldConfig, Labels } from '@grafana/data';

export enum QueryType {
  AWG = 'AWG',
  Easing = 'easing',
}

export enum WaveformType {
  Sin = 'Sin',
  Square = 'Square',
  Triangle = 'Triangle',
  Sawtooth = 'Sawtooth',
  Noise = 'Noise',
  CSV = 'CSV',
  Calculation = 'Calculation',
}

export interface WaveformArgs {
  type: WaveformType;
  period: string; // converted to seconds
  amplitude: number;
  offset: number;
  phase: number; // 0-1 shift in the period
  duty?: number; // % of the period that a squarewave is high
  points?: number[]; // for CSV
  args?: string; // Ease function for CSV, Expression for 'Calculation'
}

export interface SignalArgs {
  name: string;
  component: WaveformArgs[];
  config?: FieldConfig;
  labels?: Labels;
}

export interface SignalQuery extends DataQuery {
  queryType?: QueryType;
  signals?: SignalArgs[];
  ease?: string; // ease function matcher
}

export interface SignalDatasourceOptions extends DataSourceJsonData {
  // nothing for now
}

export interface SignalSecureJsonData {
  // nothing for now
}
