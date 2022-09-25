import { IsInt, IsString, Min, MinLength } from 'class-validator';

export class InviteToRoomDto {
  @IsString()
  @MinLength(13)
  roomId: string;

  @IsInt()
  @Min(1)
  guestId: number;
}
