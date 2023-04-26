import { NestFactory } from '@nestjs/core';
import {
  FastifyAdapter,
  NestFastifyApplication,
} from '@nestjs/platform-fastify';
import { AppModule } from './app.module';
import { AllExceptionFilter, HttpExceptionFilter } from './common/filter';
import { ValidationPipe } from '@nestjs/common';
import { AppConfig } from './config/app.config';
import { ChatAdapter } from './module/chat/chat.adapter';
import { PrismaService, RedisService } from './module/common';

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
  const chatAdapter = new ChatAdapter(app, redisService.client);
  await chatAdapter.connectToRedis();
  app.useWebSocketAdapter(chatAdapter);

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
