import { forwardRef, Module } from '@nestjs/common';
import { ChatModule } from '../chat/chat.module';
import { RoomController } from './room.controller';
import { RoomService } from './room.service';
import { PlayerModule } from '../player/player.module';
import { AuthModule } from '../auth/auth.module';

@Module({
  imports: [
    forwardRef(() => ChatModule),
    forwardRef(() => PlayerModule),
    AuthModule,
  ],
  controllers: [RoomController],
  providers: [RoomService],
  exports: [RoomService],
})
export class RoomModule {}
