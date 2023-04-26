import { IsInt, IsNotEmpty, IsPositive, IsString } from 'class-validator';
import { PlayerId } from 'src/module/user/player.type';

export class TransferOwnershipDto {
  @IsString()
  @IsNotEmpty()
  roomId: string;

  @IsInt()
  @IsPositive()
  candidateId: PlayerId;
}
