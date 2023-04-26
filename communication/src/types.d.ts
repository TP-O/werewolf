import { User } from '@prisma/client';
import { ListenEvent } from './enum/event.enum';
import { PlayerId } from './module/user/player.type';

declare module 'socket.io' {
  class Socket {
    userId?: PlayerId;
    eventName?: ListenEvent;
  }
}

declare module 'fastify' {
  export class FastifyRequest {
    user?: User;
  }
}
