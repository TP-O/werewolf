import { Module } from '@nestjs/common';
import { ConnectionService } from './chat/connection.service';
import { TextChatGateway } from './chat/text/text-chat.gateway';
import { PrismaService } from './chat/prisma.service';
import { AuthService } from './chat/auth.service';
import { UserService } from './chat/user.service';

@Module({
  imports: [],
  controllers: [],
  providers: [
    PrismaService,
    ConnectionService,
    AuthService,
    UserService,
    TextChatGateway,
  ],
})
export class AppModule {}
