import { IsNotEmpty, IsString } from 'class-validator';
import { RoomId } from '../room.type';

export class LeaveRoomDto {
  @IsString()
  @IsNotEmpty()
  id!: RoomId;
}
