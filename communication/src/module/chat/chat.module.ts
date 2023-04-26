import { Module } from '@nestjs/common';
import { PlayerModule } from '../player';
import { RoomModule, RoomService } from '../room';
import { AuthModule } from '../auth';
import { ChatService } from './chat.service';
import { ChatGateway } from '.';

@Module({
  imports: [PlayerModule, RoomModule, AuthModule],
  providers: [ChatGateway, ChatService, ChatService, RoomService],
  exports: [ChatGateway],
})
export class ChatModule {}
