import { DataSourcePlugin } from "@grafana/data";
import { DataSource } from "./DataSource";
import { ConfigEditor, QueryEditor } from "./components";
import { DroneQuery, WaveformDatasourceOptions } from "./types";

export const plugin = new DataSourcePlugin<
  DataSource,
  DroneQuery,
  WaveformDatasourceOptions
>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
