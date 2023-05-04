import { Optional } from '@nestjs/common';
import { IsString, Length } from 'class-validator';

export class CreateEmptyRoomDto {
  @IsString()
  @Length(5, 25)
  @Optional()
  password?: string;
}
