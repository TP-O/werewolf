import { IsNotEmpty, IsNumber, IsString, Max, Min } from 'class-validator';

export class RedisConfig {
  @IsString()
  @IsNotEmpty()
  public readonly host!: string;

  @IsNumber()
  @Min(1025)
  @Max(65535)
  public readonly port!: number;

  @IsString()
  @IsNotEmpty()
  public readonly password!: string;
}
