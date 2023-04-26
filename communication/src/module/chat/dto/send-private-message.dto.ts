import { IsInt, IsNotEmpty, IsPositive, IsString } from 'class-validator';
import { PlayerId } from 'src/module/player';

export class SendPrivateMessageDto {
  @IsInt()
  @IsPositive()
  receiverId: PlayerId;

  @IsString()
  @IsNotEmpty()
  content: string;
}
