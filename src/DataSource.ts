import {
  DataSourceInstanceSettings,
  DataQueryResponse,
  DataQueryRequest,
  DataFrame,
  LiveChannelAddress,
  LiveChannelScope,
} from '@grafana/data';
import { DataSourceWithBackend, getLiveMeasurementsObserver } from '@grafana/runtime';

import { SignalQuery, SignalDatasourceOptions, QueryType, SignalCustomMeta } from './types';

import { Observable, of } from 'rxjs';
import { switchMap } from 'rxjs/operators';

export class DataSource extends DataSourceWithBackend<SignalQuery, SignalDatasourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<SignalDatasourceOptions>) {
    super(instanceSettings);
  }

  query(request: DataQueryRequest<SignalQuery>): Observable<DataQueryResponse> {
    const streamQuery = request.targets.find((q) => q.queryType === QueryType.Streams);
    if (streamQuery && !streamQuery.oneshot) {
      if (request.targets.length > 1) {
        throw new Error('stream can only support one stream query at once');
      }
      return super.query(request).pipe(
        switchMap((res) => {
          if (res.data.length === 1) {
            const frame = res.data[0] as DataFrame;
            const meta = frame.meta?.custom as SignalCustomMeta;
            if (meta?.streamKey) {
              const addr: LiveChannelAddress = {
                scope: LiveChannelScope.Grafana,
                namespace: 'measurements',
                path: meta.streamKey,
              };

              console.log('streaming from:', addr, frame.fields);

              return getLiveMeasurementsObserver(addr, request.requestId);
            }
          }
          return of(res);
        })
      );
    }
    return super.query(request);
  }
}
