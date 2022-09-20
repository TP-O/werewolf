import { Module } from '@nestjs/common';
import { AuthService } from 'src/service/auth.service';
import { PrismaService } from 'src/service/prisma.service';
import { UserService } from 'src/service/user.service';
import { ConnectionService } from './service/connection.service';
import { MessageService } from './service/message.service';
import { RoomService } from './service/room.service';
import { TextChatGateway } from './text-chat.gateway';

@Module({
  imports: [],
  providers: [
    TextChatGateway,
    ConnectionService,
    MessageService,
    RoomService,
    AuthService,
    PrismaService,
    UserService,
  ],
})
export class ChatModule {}
