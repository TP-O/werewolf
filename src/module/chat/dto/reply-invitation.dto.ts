import { IsBoolean, IsString, MinLength } from 'class-validator';

export class ReplyInvitationDto {
  @IsString()
  @MinLength(13)
  roomId: string;

  @IsBoolean()
  isAccpeted: boolean;
}
