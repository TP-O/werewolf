import { Global, Module } from '@nestjs/common';
import { FirebaseService } from './service/firebase.service';
import { PrismaService } from './service/prisma.service';
import { RedisService } from './service/redis.service';
import { LoggerService } from './service/logger.service';

@Global()
@Module({
  providers: [FirebaseService, PrismaService, RedisService, LoggerService],
  exports: [FirebaseService, PrismaService, RedisService, LoggerService],
})
export class CommonModule {}
