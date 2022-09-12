import { Socket as OriginalSocket } from 'socket.io';

declare module 'socket.io' {
  export class Socket extends OriginalSocket {
    userId?: number;
  }
}
