import {
  DataSourceInstanceSettings,
  DataQueryResponse,
  DataQueryRequest,
  DataFrame,
  LiveChannelAddress,
  LiveChannelScope,
  Field,
} from '@grafana/data';
import { DataSourceWithBackend, getLiveMeasurementsObserver } from '@grafana/runtime';

import { SignalQuery, SignalDatasourceOptions, QueryType, SignalCustomMeta } from './types';

import { Observable, of } from 'rxjs';
import { switchMap, map } from 'rxjs/operators';

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
            // This frame has schema + metadata
            const frame = res.data[0] as DataFrame;
            const meta = frame.meta?.custom as SignalCustomMeta;
            if (meta?.streamKey) {
              const addr: LiveChannelAddress = {
                scope: LiveChannelScope.Grafana,
                namespace: 'measurements',
                path: meta.streamKey,
              };

              console.log('streaming from:', addr, frame.fields);

              const byName = new Map<string, Field>();
              for (const f of frame.fields) {
                byName.set(f.name, f);
              }

              // Ling the streaming data to the schema :(
              return getLiveMeasurementsObserver(addr, request.requestId).pipe(
                map((r) => {
                  const frame = r.data[0] as DataFrame;
                  if (frame) {
                    const fields = frame.fields.map((f) => {
                      const orig = byName.get(f.name);
                      if (orig) {
                        return { ...f, config: orig.config };
                      }
                      return f;
                    });

                    return {
                      ...r,
                      data: [
                        {
                          ...frame,
                          fields,
                        },
                      ],
                    };
                  }
                  return r;
                })
              );
            }
          }
          return of(res);
        })
      );
    }
    return super.query(request);
  }
}

// .pipe(
//   map((evt) => {
//     if (isLiveChannelMessageEvent(evt)) {
//       rsp.data = evt.message.getData(query);
//       if (!rsp.data.length) {
//         // ?? skip when data is empty ???
//       }
//       delete rsp.error;
//       rsp.state = LoadingState.Streaming;
//     } else if (isLiveChannelStatusEvent(evt)) {
//       if (evt.error != null) {
//         rsp.error = rsp.error;
//         rsp.state = LoadingState.Error;
//       }
//     }
//     return { ...rsp }; // send event on all status messages
//   })
// );
