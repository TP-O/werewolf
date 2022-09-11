import { Module } from '@nestjs/common';
import { ConnectionService } from './chat/connection.service';
import { TextChatGateway } from './chat/text/text-chat.gateway';
import { PrismaService } from './chat/prisma.service';
import { AuthService } from './chat/auth.service';
import { UserService } from './chat/user.service';
import { MessageService } from './chat/message.service';

@Module({
  imports: [],
  controllers: [],
  providers: [
    PrismaService,
    ConnectionService,
    AuthService,
    UserService,
    MessageService,
    TextChatGateway,
  ],
})
export class AppModule {}
