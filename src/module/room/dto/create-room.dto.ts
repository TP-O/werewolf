import { IsBoolean, IsNotEmpty } from 'class-validator';

export class CreateRoomDto {
  @IsBoolean()
  @IsNotEmpty()
  isPublic: boolean;
}
