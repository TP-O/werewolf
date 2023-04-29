import { IsNotEmpty, IsString } from 'class-validator';
import { PlayerId } from 'src/module/player';
import { RoomId } from '../room.type';

export class TransferOwnershipDto {
  @IsString()
  @IsNotEmpty()
  id!: RoomId;

  @IsString()
  @IsNotEmpty()
  newOwnerId!: PlayerId;
}
