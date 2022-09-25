import { Module } from '@nestjs/common';
import { CommunicationModule } from './module/communication/communication.module';

@Module({
  imports: [CommunicationModule],
})
export class AppModule {}
