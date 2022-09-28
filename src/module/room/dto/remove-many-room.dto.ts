import { ArrayMinSize, IsString } from 'class-validator';

export class RemoveManyRoomDto {
  @IsString({ each: true })
  @ArrayMinSize(1)
  ids: string[];
}
