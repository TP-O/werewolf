import { ArrayUnique, IsNotEmpty, IsOptional, IsString } from 'class-validator';
import { PlayerId } from 'src/module/player/player.type';
import { RoomId } from '../room.type';

export class RemoveMemberRoomsDto {
  @IsString({ each: true })
  @IsNotEmpty({ each: true })
  @IsOptional()
  @ArrayUnique()
  ids?: RoomId[];

  @IsString()
  @IsNotEmpty()
  memberId!: PlayerId;
}
