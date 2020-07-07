import { DataSourceInstanceSettings } from "@grafana/data";
import { DataSourceWithBackend } from "@grafana/runtime";

import { DroneQuery, WaveformDatasourceOptions } from "./types";

export class DataSource extends DataSourceWithBackend<
  DroneQuery,
  WaveformDatasourceOptions
> {
  constructor(
    instanceSettings: DataSourceInstanceSettings<WaveformDatasourceOptions>
  ) {
    super(instanceSettings);
  }
}
