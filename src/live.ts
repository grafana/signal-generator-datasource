import { Observable } from 'rxjs';
import { AWGQuery } from './types';
import { DataQueryResponse, CircularDataFrame, KeyValue, FieldType, LoadingState } from '@grafana/data';
import * as minimatch from 'minimatch';
import { ObservableSubject } from './subject';
import { doConnect } from 'broker/streams';

interface InfluxMessage {
  name: string;
  fields: Record<string, number>;
  tags: Record<string, string>;
  timestamp: number;
}

interface InfluxMessageCollection {
  metrics: InfluxMessage[];
}

let socket: any = undefined;
let subjects: ObservableSubject[] = [];
let lastMessage = 0;

export function listenToSocket(request: AWGQuery): Observable<DataQueryResponse> {
  if (!socket) {
    console.log('Connecting to websocket...');
    socket = doConnect(processMsg);
  }

  const subject = new ObservableSubject(request);
  subjects.push(subject);
  return subject.asObservable();
}

function nameGlobMatches(msg: InfluxMessage, filters: KeyValue<string>): boolean {
  for (let key of Object.keys(filters)) {
    let value = filters[key];
    if (key === 'namepass' && minimatch.match([msg.name], value).length === 0) {
      return false;
    }

    if (key === 'namedrop' && minimatch.match([msg.name], value).length !== 0) {
      return false;
    }
  }
  return true;
}

function filterFields(fields: Record<string, number>, filters: KeyValue<string>): Record<string, number> {
  const pass = filters['fieldpass'];
  const drop = filters['fielddrop'];

  if (!pass && !drop) {
    return fields;
  }

  return Object.keys(fields).reduce((f, key) => {
    if (pass && minimatch.match([key], pass).length !== 0) {
      f[key] = fields[key];
    }
    if (drop && minimatch.match([key], drop).length === 0) {
      f[key] = fields[key];
    }
    return f;
  }, {} as Record<string, number>);
}

function filterTags(tags: Record<string, string>, filters: KeyValue<string>): Record<string, string> {
  return {};
}

function processMsg(msgs: InfluxMessageCollection) {
  subjects.forEach(s => {
    const subject = s.subject;
    const data = s.data;
    //    const maxdata: number = msgs.metrics.length - (parseInt(s.filters['maxdata'], 10) || msgs.metrics.length);
    msgs.metrics
      .filter(msg => nameGlobMatches(msg, s.filters))
      .forEach((msg, i) => {
        const name = msg.name as string;
        let df = data[name];
        if (!df) {
          df = new CircularDataFrame({
            append: 'tail',
            capacity: 500, //this.query.buffer,
          });
          df.name = name;
          df.addField({ name: 'timestamp', type: FieldType.time }, 0);
          data[name] = df;
        }

        // if (i > maxdata) {
        const row = {
          timestamp: msg.timestamp, // millis, not seconds * 1000,
          ...filterFields(msg.fields, s.filters),
          ...filterTags(msg.tags, s.filters),
        };

        df.add(row, true);
        // }

        const elapsed = Date.now() - lastMessage;
        if (elapsed > 0) {
          subject.next({
            data: Object.values(data),
            key: 'ws',
            state: LoadingState.Done,
          });
          lastMessage = Date.now();
        }
      });
  });
}
