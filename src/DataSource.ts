import {
  DataSourceInstanceSettings,
  DataQueryResponse,
  DataQueryRequest,
  DataFrame,
  LiveChannelAddress,
  LiveChannelScope,
  Field,
} from '@grafana/data';
import { DataSourceWithBackend, getLiveDataStream } from '@grafana/runtime';

import { SignalQuery, SignalDatasourceOptions, QueryType, SignalCustomMeta } from './types';

import { Observable, of } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { LiveMeasurementsSupport } from 'support';

export class DataSource extends DataSourceWithBackend<SignalQuery, SignalDatasourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<SignalDatasourceOptions>) {
    super(instanceSettings);

    // Channels managed by this datasource instance
    this.channelSupport = new LiveMeasurementsSupport();
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
            // This frame has schema + metadata
            const frame = res.data[0] as DataFrame;
            const meta = frame.meta?.custom as SignalCustomMeta;
            if (meta?.streamKey) {
              const addr: LiveChannelAddress = {
                scope: LiveChannelScope.DataSource,
                namespace: this.uid,
                path: meta.streamKey,
              };

              console.log('streaming from:', addr, frame.fields);

              const byName = new Map<string, Field>();
              for (const f of frame.fields) {
                byName.set(f.name, f);
              }

              // Ling the streaming data to the schema :(
              return getLiveDataStream({ addr, key: request.requestId });
            }
          }
          return of(res);
        })
      );
    }
    return super.query(request);
  }
}
