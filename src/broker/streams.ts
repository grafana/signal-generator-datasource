import Centrifuge, { PublicationContext } from 'centrifuge/dist/centrifuge.protobuf';
import SockJS from 'sockjs-client';

const centrifuge = new Centrifuge('http://localhost:3007/broker/sockjs', {
  debug: true,
  sockjs: SockJS,
});

centrifuge.setToken('ABCD');

centrifuge.on('connect', function(context) {
  console.log('CONNECT', context);
});

centrifuge.on('disconnect', function(context) {
  console.log('disconnect', context);
});

// Server side function
centrifuge.on('publish', function(ctx) {
  console.log('Publication from server-side channel', ctx);
});

export function doConnect(onMsg: (msgs: any) => void) {
  centrifuge.connect();
  return centrifuge.subscribe('example', (message: PublicationContext) => {
    onMsg(message.data);
  });
}
