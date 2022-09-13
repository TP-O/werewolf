import { Module } from '@nestjs/common';
import { CommonModule } from '../common/common.module';
import { ConnectionService } from './connection.service';
import { TextChatGateway } from './/text-chat.gateway';
import { MessageService } from './message.service';
import { RoomService } from './room.service';

@Module({
  imports: [CommonModule],
  providers: [TextChatGateway, ConnectionService, MessageService, RoomService],
})
export class ChatModule {}
