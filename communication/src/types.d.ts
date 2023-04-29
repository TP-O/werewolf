import { Player } from '@prisma/client';
import { ListenEvent } from './module/chat/chat.enum';
import { OnlinePlayer, PlayerId } from './module/player/player.type';

declare module 'socket.io' {
  class Socket {
    playerId?: PlayerId;
    eventName?: ListenEvent;
  }
}

declare module 'fastify' {
  export class FastifyRequest {
    player?: Player | OnlinePlayer;
  }
}
