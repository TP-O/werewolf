import {
  ArrayUnique,
  IsNotEmpty,
  IsNumber,
  IsString,
  Min,
} from 'class-validator';
import { PlayerId } from 'src/module/user/player.type';

export class AddToRoomDto {
  @IsString()
  @IsNotEmpty()
  roomId: string;

  @IsNumber({}, { each: true })
  @Min(1, { each: true })
  @ArrayUnique()
  memberIds: PlayerId[];
}
