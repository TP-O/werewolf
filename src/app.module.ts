import { Module } from '@nestjs/common';
import { ChatModule } from './module/chat/chat.module';

@Module({
  imports: [ChatModule],
})
export class AppModule {}
