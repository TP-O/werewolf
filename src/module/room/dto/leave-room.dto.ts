import { IsString, MinLength } from 'class-validator';

export class LeaveRoomDto {
  @IsString()
  @MinLength(13)
  roomId: string;
}
