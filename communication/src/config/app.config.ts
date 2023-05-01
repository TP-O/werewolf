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
import { AppEnv } from 'src/common/enum';

class CorsOptions {
  @IsArray()
  @ArrayMinSize(1)
  @IsString({ each: true })
  public readonly origins!: string[];
}

export class AppConfig {
  @IsIn(Object.values(AppEnv))
  public readonly env!: AppEnv;

  @IsBoolean()
  readonly debug!: boolean;

  @IsNumber()
  @Min(10)
  @Max(65535)
  public readonly port!: number;

  @IsString()
  @MinLength(20)
  public readonly secret!: string;

  @Type(() => CorsOptions)
  @ValidateNested()
  public readonly cors!: CorsOptions;
}
