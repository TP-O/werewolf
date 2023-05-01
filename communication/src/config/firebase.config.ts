import { Type } from 'class-transformer';
import {
  IsNotEmpty,
  IsOptional,
  IsString,
  ValidateNested,
} from 'class-validator';

class FirebaseClientConfig {
  @IsString()
  @IsOptional()
  public readonly apiKey!: string;

  @IsString()
  @IsOptional()
  public readonly authDomain!: string;

  @IsString()
  @IsOptional()
  public readonly projectId!: string;

  @IsString()
  @IsOptional()
  public readonly storageBucket!: string;

  @IsString()
  @IsOptional()
  public readonly messagingSenderId!: string;

  @IsString()
  @IsOptional()
  public readonly appId!: string;

  @IsString()
  @IsOptional()
  public readonly measurementId!: string;
}

export class FirebaseConfig {
  @IsString()
  @IsNotEmpty()
  public readonly productId!: string;

  @IsString()
  @IsNotEmpty()
  public readonly privateKey!: string;

  @IsString()
  @IsNotEmpty()
  public readonly email!: string;

  /**
   * Client firebase config is required in development env only.
   */
  @Type(() => FirebaseClientConfig)
  @ValidateNested()
  public readonly client!: FirebaseClientConfig;
}
