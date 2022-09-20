import { IsBoolean } from 'class-validator';

export class CreateRoomDto {
  @IsBoolean()
  isPublic: boolean;
}
