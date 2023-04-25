import { NestFactory } from '@nestjs/core';
import {
  FastifyAdapter,
  NestFastifyApplication,
} from '@nestjs/platform-fastify';
import { RedisIoAdapter } from './adapter';
import { AppModule } from './app.module';
import { AllExceptionFilter, HttpExceptionFilter } from './common/filter';
import { PrismaService } from './common/service/prisma.service';
import { ValidationPipe } from '@nestjs/common';
import { RedisService } from './common/service/redis.service';
import { AppConfig } from './config/app';

async function bootstrap() {
  const app = await NestFactory.create<NestFastifyApplication>(
    AppModule,
    new FastifyAdapter(),
  );
  const config = app.get(AppConfig);

  app.enableCors({
    origin: config.cors.origins,
    methods: ['GET', 'POST'],
    credentials: true,
  });

  const redisService = app.get(RedisService);
  const redisIoAdapter = new RedisIoAdapter(app, redisService.client);
  await redisIoAdapter.connectToRedis();
  app.useWebSocketAdapter(redisIoAdapter);

  const prismaService = app.get(PrismaService);
  await prismaService.enableShutdownHooks(app);

  app.setGlobalPrefix('api/v1');
  app.useGlobalFilters(new AllExceptionFilter(), new HttpExceptionFilter());
  app.useGlobalPipes(
    new ValidationPipe({
      whitelist: true,
      stopAtFirstError: false,
      transform: true,
    }),
  );

  await app.listen(config.port, '0.0.0.0');
}

bootstrap();
