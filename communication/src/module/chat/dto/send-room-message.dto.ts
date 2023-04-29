import { IsNotEmpty, IsString } from 'class-validator';

export class SendRoomMessageDto {
  @IsString()
  @IsNotEmpty()
  roomId!: string;

  @IsString()
  @IsNotEmpty()
  content!: string;
}
