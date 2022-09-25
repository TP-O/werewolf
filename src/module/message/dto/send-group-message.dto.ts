import { IsNotEmpty, IsString, MinLength } from 'class-validator';

export class SendGroupMessageDto {
  @IsString()
  @MinLength(13)
  roomId: string;

  @IsString()
  @IsNotEmpty()
  content: string;
}
