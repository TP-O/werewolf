import { IsNotEmpty, IsString } from 'class-validator';

export class FirebaseConfig {
  @IsString()
  @IsNotEmpty()
  public readonly productId!: string;

  @IsString()
  @IsNotEmpty()
  public readonly privateKey!: string;

  @IsString()
  @IsNotEmpty()
  public readonly clientEmail!: string;
}
