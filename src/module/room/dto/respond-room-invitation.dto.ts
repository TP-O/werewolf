import { IsBoolean, IsString, MinLength } from 'class-validator';

export class RespondRoomInvitationDto {
  @IsString()
  @MinLength(13)
  roomId: string;

  @IsBoolean()
  isAccpeted: boolean;
}
