import { IoAdapter } from '@nestjs/platform-socket.io';
import Redis from 'ioredis';
import { ServerOptions } from 'socket.io';
import { createAdapter } from '@socket.io/redis-adapter';
import { NestFastifyApplication } from '@nestjs/platform-fastify';

export class ChatAdapter extends IoAdapter {
  private adapterConstructor!: ReturnType<typeof createAdapter>;

  constructor(private readonly redis: Redis, app: NestFastifyApplication) {
    super(app);
  }

  async connectToRedis(): Promise<void> {
    this.adapterConstructor = createAdapter(this.redis, this.redis.duplicate());
  }

  createIOServer(port: number, options?: ServerOptions) {
    const server = super.createIOServer(port, options);
    server.adapter(this.adapterConstructor);
    return server;
  }
}
