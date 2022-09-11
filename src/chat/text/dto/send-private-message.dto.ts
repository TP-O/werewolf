import { IsInt, IsNotEmpty, IsString, Min } from 'class-validator';

export class SendPrivateMessageDto {
  @IsInt()
  @Min(1)
  receivedId: number;

  @IsString()
  @IsNotEmpty()
  content: string;
}
