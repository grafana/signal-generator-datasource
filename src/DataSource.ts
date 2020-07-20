import { DataSourceInstanceSettings, DataQueryRequest, DataQueryResponse, KeyValue } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { Observable } from 'rxjs';

import { AWGQuery, AWGQueryType, AWGDatasourceOptions } from './types';
import { listenToSocket } from 'live';

export class DataSource extends DataSourceWithBackend<AWGQuery, AWGDatasourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<AWGDatasourceOptions>) {
    super(instanceSettings);
  }

  query(request: DataQueryRequest<AWGQuery>): Observable<DataQueryResponse> {
    for (const target of request.targets) {
      if (target.queryType === AWGQueryType.Stream) {
        console.log('TODO... open websocket!');

        return listenToSocket(target);
      }
    }
    return super.query(request);
  }
}
