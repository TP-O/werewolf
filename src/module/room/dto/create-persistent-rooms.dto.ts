import { Type } from 'class-transformer';
import {
  ArrayMinSize,
  ArrayUnique,
  IsArray,
  IsNotEmpty,
  IsNumber,
  IsString,
  Min,
  ValidateNested,
} from 'class-validator';
import { Room } from '../room.type';

class CreatePersistentRoomDto implements Room {
  @IsString()
  @IsNotEmpty()
  id: string;

  @IsNumber({}, { each: true })
  @Min(1, { each: true })
  @ArrayUnique()
  memberIds: number[];

  ownerId: number;

  gameId: number;

  isPublic: boolean;

  isPersistent = true;

  waitingIds: number[] = [];

  refusedIds: number[] = [];
}

export class CreatePersistentRoomsDto {
  @IsNumber()
  @Min(1)
  gameId: number;

  @IsArray()
  @ArrayMinSize(1)
  @ArrayUnique((value: CreatePersistentRoomDto) => value.id)
  @ValidateNested({ each: true })
  @Type(() => CreatePersistentRoomDto)
  rooms: CreatePersistentRoomDto[];
}
