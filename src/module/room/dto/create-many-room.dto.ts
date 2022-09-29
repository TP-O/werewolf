import { Type } from 'class-transformer';
import {
  ArrayMinSize,
  ArrayUnique,
  Equals,
  IsArray,
  IsBoolean,
  IsNotEmpty,
  IsNumber,
  IsString,
  Min,
  ValidateIf,
  ValidateNested,
} from 'class-validator';

class CreateRoomDto {
  @IsString()
  @IsNotEmpty()
  id: string;

  @IsBoolean()
  isPublic: boolean;

  @IsBoolean()
  isPersistent: boolean;

  @ValidateIf((dto: CreateRoomDto, v) => !dto.memberIds?.includes(v))
  @Equals(undefined, { message: '$property must be contained in memberIds' })
  ownerId: number;

  @IsNumber({}, { each: true })
  @Min(1, { each: true })
  @ArrayUnique()
  memberIds: number[];
}

export class CreateManyRoomDto {
  @IsArray()
  @ArrayMinSize(1)
  @ArrayUnique((value: CreateRoomDto) => value.id)
  @ValidateNested({ each: true })
  @Type(() => CreateRoomDto)
  rooms: CreateRoomDto[];
}
