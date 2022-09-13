import { IsString, MinLength } from 'class-validator';

export class JoinRoomDto {
  @IsString()
  @MinLength(13)
  id: string;
}
