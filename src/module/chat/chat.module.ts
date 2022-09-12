import { Module } from '@nestjs/common';
import { CommonModule } from '../common/common.module';
import { ConnectionService } from './connection.service';
import { TextChatGateway } from './/text-chat.gateway';
import { MessageService } from './message.service';

@Module({
  imports: [CommonModule],
  providers: [TextChatGateway, ConnectionService, MessageService],
})
export class ChatModule {}
