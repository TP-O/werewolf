import { IsInt, IsNotEmpty, IsString, Min } from 'class-validator';
import { Exclude } from 'class-transformer';

export class SendPrivateMessageDto {
  @IsInt()
  @Min(1)
  @Exclude({ toPlainOnly: true })
  receiverId: number;

  @IsString()
  @IsNotEmpty()
  content: string;
}
