import { Module } from '@nestjs/common';
import { ChatModule } from './module/chat/chat.module';
import { RoomModule } from './module/room/room.module';
import { PlayerModule } from './module/player/player.module';
import { TypedConfigModule } from 'nest-typed-config';
import { RootConfig } from './config/main';
import { CommonModule } from './module/common/global.module';
import { AuthModule } from './module/auth/auth.module';
import { loadConfig } from './utils/load-config';

@Module({
  imports: [
    CommonModule,
    AuthModule,
    ChatModule,
    PlayerModule,
    RoomModule,
    TypedConfigModule.forRoot({
      schema: RootConfig,
      load: loadConfig,
    }),
  ],
  providers: [],
})
export class AppModule {}
