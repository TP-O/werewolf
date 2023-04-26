import { IsInt, IsNotEmpty, IsPositive, IsString } from 'class-validator';
import { PlayerId } from 'src/module/user/player.type';

export class InviteToRoomDto {
  @IsString()
  @IsNotEmpty()
  roomId: string;

  @IsInt()
  @IsPositive()
  guestId: PlayerId;
}
