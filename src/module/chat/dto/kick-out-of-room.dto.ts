import { IsInt, IsString, Min, MinLength } from 'class-validator';

export class KickOutOfRoomDto {
  @IsString()
  @MinLength(13)
  roomId: string;

  @IsInt()
  @Min(1)
  memberId: number;
}
