import { Type } from 'class-transformer';
import {
  ArrayMinSize,
  ArrayUnique,
  Equals,
  IsArray,
  IsBoolean,
  IsNumber,
  Min,
  ValidateIf,
  ValidateNested,
} from 'class-validator';
import { Room } from '../room.type';
import { PlayerId } from 'src/module/user/player.type';

class CreateRoomDto implements Room {
  @IsBoolean()
  isPublic: boolean;

  @ValidateIf((dto: CreateRoomDto, v) => !dto.memberIds?.includes(v))
  @Equals(undefined, { message: '$property must be contained in memberIds' })
  ownerId: PlayerId;

  @IsNumber({}, { each: true })
  @Min(1, { each: true })
  @ArrayUnique()
  memberIds: PlayerId[];

  id: string;

  isPersistent = false;

  isMuted = false;

  gameId = 0;

  waitingIds: PlayerId[] = [];

  refusedIds: PlayerId[] = [];
}

export class CreateTemporaryRoomsDto {
  @IsArray()
  @ArrayMinSize(1)
  @ValidateNested({ each: true })
  @Type(() => CreateRoomDto)
  rooms: CreateRoomDto[];
}
