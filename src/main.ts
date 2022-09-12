import { NestFactory, Reflector } from '@nestjs/core';
import {
  FastifyAdapter,
  NestFastifyApplication,
} from '@nestjs/platform-fastify';
import { RedisIoAdapter } from './adapter/redis.adapter';
import { AppModule } from './app.module';
import { AppConfig } from './config/app.config';
import { RolesGuard } from './guard/roles.guard';
import { PrismaService } from './module/common/prisma.service';

async function bootstrap() {
  const app = await NestFactory.create<NestFastifyApplication>(
    AppModule,
    new FastifyAdapter(),
  );

  const redisIoAdapter = new RedisIoAdapter(app);
  await redisIoAdapter.connectToRedis();
  app.useWebSocketAdapter(redisIoAdapter);

  const prismaService = app.get(PrismaService);
  await prismaService.enableShutdownHooks(app);

  app.useGlobalGuards(new RolesGuard(app.get(Reflector)));

  await app.listen(AppConfig.port, '0.0.0.0');
}

bootstrap();
