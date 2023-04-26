import { Global, Module } from '@nestjs/common';
import { FirebaseService } from './service/firebase.service';
import { PrismaService } from './service/prisma.service';
import { RedisService } from './service/redis.service';

@Global()
@Module({
  providers: [FirebaseService, PrismaService, RedisService],
  exports: [FirebaseService, PrismaService, RedisService],
})
export class CommonModule {
  //
}
