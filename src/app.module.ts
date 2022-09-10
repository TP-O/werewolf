import { CacheModule, Module } from '@nestjs/common';
import { ConnectionService } from './chat/connection.service';
import { TextChatGateway } from './chat/text/text-chat.gateway';
import { CacheConfig } from './config/cache.config';
import { PrismaService } from './chat/prisma.service';
import { AuthService } from './chat/auth.service';

@Module({
  imports: [CacheModule.register(CacheConfig)],
  controllers: [],
  providers: [PrismaService, ConnectionService, AuthService, TextChatGateway],
})
export class AppModule {}
