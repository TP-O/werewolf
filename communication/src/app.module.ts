import { Module } from '@nestjs/common';
import { ChatModule } from './module/chat/chat.module';
import { RoomModule } from './module/room/room.module';
import { PlayerModule } from './module/player/player.module';
import { TypedConfigModule, fileLoader } from 'nest-typed-config';
import { RootConfig } from './config/main';
import { CommonModule } from './module/common/global.module';
import { AuthModule } from './module/auth/auth.module';

@Module({
  imports: [
    CommonModule,
    AuthModule,
    ChatModule,
    PlayerModule,
    RoomModule,
    TypedConfigModule.forRoot({
      schema: RootConfig,
      load: fileLoader({
        absolutePath: process.env.CONFIG_FILE ?? 'config.yaml',
      }),
    }),
  ],
  providers: [],
})
export class AppModule {}
