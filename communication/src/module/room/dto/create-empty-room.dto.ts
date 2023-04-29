import { IsString, Length } from 'class-validator';

export class CreateEmptyRoomDto {
  @IsString()
  @Length(5, 25)
  password?: string;
}
