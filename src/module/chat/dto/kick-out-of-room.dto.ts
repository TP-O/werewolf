import { IsInt, IsString, Min, MinLength } from 'class-validator';

export class KickOutOfRoomDto {
  @IsString()
  @MinLength(13)
  id: string;

  @IsInt()
  @Min(1)
  memberId: number;
}
