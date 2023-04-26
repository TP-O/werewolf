import { IsInt, IsNotEmpty, IsPositive, IsString } from 'class-validator';
import { PlayerId } from 'src/module/user/player.type';

export class SendPrivateMessageDto {
  @IsInt()
  @IsPositive()
  receiverId: PlayerId;

  @IsString()
  @IsNotEmpty()
  content: string;
}
