import { DataSourceInstanceSettings, DataQueryRequest, DataQueryResponse } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { Observable } from 'rxjs';

import { AWGQuery, AWGQueryType, AWGDatasourceOptions } from './types';
import { listenToSocket } from 'live';
import { checkConnection } from 'broker/streams';

export class DataSource extends DataSourceWithBackend<AWGQuery, AWGDatasourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<AWGDatasourceOptions>) {
    super(instanceSettings);
  }

  query(request: DataQueryRequest<AWGQuery>): Observable<DataQueryResponse> {
    checkConnection();

    for (const target of request.targets) {
      if (target.queryType === AWGQueryType.Stream) {
        console.log('TODO... open websocket!');

        return listenToSocket();
      }
    }
    return super.query(request);
  }
}
