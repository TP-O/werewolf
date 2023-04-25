import { Module } from '@nestjs/common';
import { CommunicationModule } from './module/communication/communication.module';
import { MessageModule } from './module/message/message.module';
import { RoomModule } from './module/room/room.module';
import { UserModule } from './module/user/user.module';
import { TypedConfigModule, fileLoader } from 'nest-typed-config';
import { RootConfig } from './config/main';

@Module({
  imports: [
    CommunicationModule,
    UserModule,
    RoomModule,
    MessageModule,
    TypedConfigModule.forRoot({
      schema: RootConfig,
      load: fileLoader({
        absolutePath: 'config.yaml',
      }),
    }),
  ],
  providers: [],
})
export class AppModule {}
