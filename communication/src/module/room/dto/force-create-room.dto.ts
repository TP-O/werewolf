import {
  ArrayMinSize,
  ArrayUnique,
  IsBoolean,
  IsNotEmpty,
  IsString,
  Length,
  ValidateNested,
} from 'class-validator';
import { PlayerId } from 'src/module/player';
import { RoomId } from '../room.type';

export class ForceCreateRoomDto {
  @IsString()
  @IsNotEmpty()
  id!: RoomId;

  @IsString()
  @IsNotEmpty()
  ownerId?: PlayerId;

  @IsString()
  @Length(5, 25)
  password?: string;

  @IsBoolean()
  isMuted!: boolean;

  @IsString({ each: true })
  @IsNotEmpty({ each: true })
  @ArrayUnique()
  @ArrayMinSize(1)
  memberIds!: PlayerId[];
}

export class ForceCreateRoomsDto {
  @ArrayMinSize(1)
  @ValidateNested({ each: true })
  rooms!: ForceCreateRoomDto[];
}
