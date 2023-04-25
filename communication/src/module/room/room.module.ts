import { forwardRef, Module } from '@nestjs/common';
import { PrismaService } from 'src/common/service/prisma.service';
import { CommunicationModule } from '../communication/communication.module';
import { RoomController } from './room.controller';
import { RoomService } from './room.service';
import { RedisService } from 'src/common/service/redis.service';

@Module({
  imports: [forwardRef(() => CommunicationModule)],
  controllers: [RoomController],
  providers: [RoomService, PrismaService, RedisService],
  exports: [RoomService],
})
export class RoomModule {
  //
}
