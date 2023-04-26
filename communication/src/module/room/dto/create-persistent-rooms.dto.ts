import { Type } from 'class-transformer';
import {
  ArrayMinSize,
  ArrayUnique,
  IsArray,
  IsBoolean,
  IsNotEmpty,
  IsNumber,
  IsPositive,
  IsString,
  ValidateNested,
} from 'class-validator';
import { Room } from '../room.type';
import { PlayerId } from 'src/module/player';

class CreatePersistentRoomDto implements Room {
  @IsString()
  @IsNotEmpty()
  id: string;

  @IsBoolean()
  isMuted: boolean;

  @IsNumber({}, { each: true })
  @IsPositive({ each: true })
  @ArrayUnique()
  memberIds: PlayerId[];

  ownerId: PlayerId;

  gameId: number;

  isPublic: boolean;

  isPersistent = true;

  waitingIds: PlayerId[] = [];

  refusedIds: PlayerId[] = [];
}

export class CreatePersistentRoomsDto {
  @IsNumber()
  @IsPositive()
  gameId: number;

  @IsArray()
  @ArrayMinSize(1)
  @ArrayUnique((value: CreatePersistentRoomDto) => value.id)
  @ValidateNested({ each: true })
  @Type(() => CreatePersistentRoomDto)
  rooms: CreatePersistentRoomDto[];
}
