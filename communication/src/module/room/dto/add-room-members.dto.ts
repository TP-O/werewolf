import {
  ArrayMinSize,
  ArrayUnique,
  IsNotEmpty,
  IsString,
} from 'class-validator';
import { PlayerId } from 'src/module/player';
import { RoomId } from '../room.type';

export class AddRoomMembersDto {
  @IsString()
  @IsNotEmpty()
  id!: RoomId;

  @IsString({ each: true })
  @IsNotEmpty({ each: true })
  @ArrayMinSize(1)
  @ArrayUnique()
  memberIds!: PlayerId[];
}
