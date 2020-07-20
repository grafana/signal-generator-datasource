import Centrifuge, { PublicationContext } from 'centrifuge/dist/centrifuge.protobuf';
import SockJS from 'sockjs-client';

const centrifuge = new Centrifuge('http://localhost:3007/broker/sockjs', {
  debug: true,
  sockjs: SockJS,
});

centrifuge.setToken('ABCD');

centrifuge.on('connect', function(context) {
  // now client connected to Centrifugo and authorized
  console.log('CONNECT', context);
});

centrifuge.on('disconnect', function(context) {
  console.log('disconnect', context);
});

// Server side function
centrifuge.on('publish', function(ctx) {
  console.log('Publication from server-side channel', ctx);
});

centrifuge.connect();

centrifuge.subscribe('simple', (message: PublicationContext) => {
  console.log('MESSAGE', message);
});

export function checkConnection() {
  console.log('BEFORE');
  console.log('IS Connected', centrifuge.isConnected());
}
