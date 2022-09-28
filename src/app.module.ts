import { Module } from '@nestjs/common';
import { CommunicationModule } from './module/communication/communication.module';
import { RoomModule } from './module/room/room.module';
import { UserModule } from './module/user/user.module';

@Module({
  imports: [CommunicationModule, UserModule, RoomModule],
})
export class AppModule {}
