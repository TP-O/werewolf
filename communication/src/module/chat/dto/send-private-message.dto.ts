import { IsNotEmpty, IsString } from 'class-validator';
import { PlayerId } from 'src/module/player/player.type';

export class SendPrivateMessageDto {
  @IsString()
  @IsNotEmpty()
  receiverId!: PlayerId;

  @IsString()
  @IsNotEmpty()
  content!: string;
}
