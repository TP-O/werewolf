import { IsInt, IsString, IsNotEmpty, IsPositive } from 'class-validator';
import { PlayerId } from 'src/module/player';

export class KickOutOfRoomDto {
  @IsString()
  @IsNotEmpty()
  roomId: string;

  @IsInt()
  @IsPositive()
  memberId: PlayerId;
}
