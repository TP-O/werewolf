import { forwardRef, Module } from '@nestjs/common';
import { CommunicationModule } from '../communication/communication.module';
import { RoomController } from './room.controller';
import { RoomService } from './room.service';

@Module({
  imports: [forwardRef(() => CommunicationModule)],
  controllers: [RoomController],
  providers: [RoomService],
  exports: [RoomService],
})
export class RoomModule {
  //
}
