import { User } from '@prisma/client';
import { FastifyRequest as OriginalFastifyRequest } from 'fastify';
import { Socket as OriginalSocket } from 'socket.io';

declare module 'socket.io' {
  export class Socket extends OriginalSocket {
    userId?: number;
  }
}

declare module 'fastify' {
  export class FastifyRequest extends OriginalFastifyRequest {
    user?: User;
  }
}
