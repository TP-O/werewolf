import { ArrayMinSize, IsString } from 'class-validator';

export class RemoveRoomsDto {
  @IsString({ each: true })
  @ArrayMinSize(1)
  ids: string[];
}
