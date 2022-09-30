import { IsBoolean, IsString, MinLength } from 'class-validator';

export class MuteRoomDto {
  @IsString()
  @MinLength(13)
  roomId: string;

  @IsBoolean()
  mute: boolean;
}
