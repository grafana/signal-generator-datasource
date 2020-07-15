import { ReplaySubject, Observable } from 'rxjs';
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';
import { DataQueryResponse, CircularDataFrame, KeyValue, FieldType, LoadingState } from '@grafana/data';

interface InfluxMessage {
  name: string;
  fields: Record<string, number>;
  tags: Record<string, string>;
  timestamp: number;
}

let socket: WebSocketSubject<InfluxMessage> | undefined;
const subject = new ReplaySubject<DataQueryResponse>(1);
const data: KeyValue<CircularDataFrame> = {};
let lastMessage = 0;

export function listenToSocket(): Observable<DataQueryResponse> {
  if (!socket) {
    console.log('Connecting to websocket...');
    socket = webSocket('ws://localhost:3003/subscribe');
    socket.subscribe(
      processMsg, // Called whenever there is a message from the server.
      err => console.log('ERROR', err), // Called if at any point WebSocket API signals some kind of error.
      () => console.log('complete') // Called when connection is closed (for whatever reason).
    );
  }
  return subject.asObservable();
}

function processMsg(msg: InfluxMessage) {
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

  const row = {
    timestamp: msg.timestamp, // millis, not seconds * 1000,
    ...msg.fields,
    ...msg.tags,
  };

  df.add(row, true);

  const elapsed = Date.now() - lastMessage;
  if (elapsed > 0) {
    subject.next({
      data: Object.values(data),
      key: 'ws',
      state: LoadingState.Done,
    });
    lastMessage = Date.now();
  }
}
