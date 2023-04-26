import { forwardRef, Module } from '@nestjs/common';
import { ChatModule } from '../chat/chat.module';
import { RoomController } from './room.controller';
import { RoomService } from './room.service';

@Module({
  imports: [forwardRef(() => ChatModule)],
  controllers: [RoomController],
  providers: [RoomService],
  exports: [RoomService],
})
export class RoomModule {
  //
}
