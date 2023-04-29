import { IsBoolean, IsNotEmpty, IsString } from 'class-validator';
import { RoomId } from '../room.type';

export class MuteRoomDto {
  @IsString()
  @IsNotEmpty()
  id!: RoomId;

  @IsBoolean()
  isMuted!: boolean;
}
