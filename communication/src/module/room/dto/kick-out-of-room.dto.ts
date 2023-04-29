import { IsString, IsNotEmpty } from 'class-validator';
import { PlayerId } from 'src/module/player';
import { RoomId } from '../room.type';

export class KickOutOfRoomDto {
  @IsString()
  @IsNotEmpty()
  id!: RoomId;

  @IsString()
  @IsNotEmpty()
  memberId!: PlayerId;
}
