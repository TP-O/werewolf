import { Module } from '@nestjs/common';
import { PlayerController } from './player.controller';
import { PlayerService } from './player.service';
import { RoomModule } from '../room';
import { AuthModule } from '../auth';

@Module({
  imports: [RoomModule, AuthModule],
  controllers: [PlayerController],
  providers: [PlayerService],
  exports: [PlayerService],
})
export class PlayerModule {}
