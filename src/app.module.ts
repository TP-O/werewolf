import { Module } from '@nestjs/common';
import { CommunicationModule } from './module/communication/communication.module';
import { UserModule } from './module/user/user.module';

@Module({
  imports: [CommunicationModule, UserModule],
})
export class AppModule {}
