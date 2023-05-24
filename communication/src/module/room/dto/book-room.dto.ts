import { IsOptional, IsString, Length } from 'class-validator';

export class BookRoomDto {
  @IsString()
  @Length(5, 25)
  @IsOptional()
  password?: string;
}
