import { IsNotEmpty, IsOptional, IsString, Length } from 'class-validator';
import { RoomId } from '../room.type';

export class JoinRoomDto {
  @IsString()
  @IsNotEmpty()
  id!: RoomId;

  @IsString()
  @Length(5, 25)
  @IsOptional()
  password?: string;
}
