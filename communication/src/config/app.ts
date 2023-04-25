import { Type } from 'class-transformer';
import {
  ArrayMinSize,
  IsArray,
  IsBoolean,
  IsIn,
  IsNumber,
  IsString,
  Max,
  Min,
  MinLength,
  ValidateNested,
} from 'class-validator';

class CorsOptions {
  @IsArray()
  @ArrayMinSize(1)
  @IsString({ each: true })
  public readonly origins!: string[];
}

export class AppConfig {
  @IsIn(['development', 'production'])
  public readonly env!: 'development' | 'production';

  @IsBoolean()
  readonly debug!: boolean;

  @IsNumber()
  @Min(1025)
  @Max(65535)
  public readonly port!: number;

  @IsString()
  @MinLength(20)
  public readonly secret!: string;

  @IsString()
  @MinLength(20)
  public readonly decrypted!: string;

  @Type(() => CorsOptions)
  @ValidateNested()
  public readonly cors!: CorsOptions;
}
