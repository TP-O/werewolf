import { IsNotEmpty, IsString, Length } from 'class-validator';
import { RoomId } from '../room.type';

export class JoinRoomDto {
  @IsString()
  @IsNotEmpty()
  id!: RoomId;

  @IsString()
  @Length(5, 25)
  password?: string;
}
