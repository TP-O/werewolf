import { IsInt, IsString, Min, MinLength } from 'class-validator';

export class TransferOwnershipDto {
  @IsString()
  @MinLength(13)
  roomId: string;

  @IsInt()
  @Min(1)
  candidateId: number;
}
