import { Module } from '@nestjs/common';
import { PlayerModule } from '../player/player.module';
import { RoomModule } from '../room/room.module';
import { AuthModule } from '../auth/auth.module';
import { ChatService } from './chat.service';
import { ChatGateway } from './chat.gateway';
import { RoomService } from '../room/room.service';

@Module({
  imports: [PlayerModule, RoomModule, AuthModule],
  providers: [ChatGateway, ChatService, ChatService, RoomService],
  exports: [ChatGateway],
})
export class ChatModule {}
