import { IsString, MinLength } from 'class-validator';

export class LeaveRoomDto {
  @IsString()
  @MinLength(13)
  id: string;
}
