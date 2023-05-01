import { Body, Controller, HttpStatus, Post, Res } from '@nestjs/common';
import { SignInDto } from './dto';
import { FastifyReply } from 'fastify';
import { initializeApp } from 'firebase/app';
import { getAuth, signInWithEmailAndPassword } from 'firebase/auth';
import { AppConfig, FirebaseConfig } from 'src/config';
import { AppEnv } from 'src/common/enum';

@Controller('auth')
export class AuthController {
  private readonly _isDev: boolean;

  constructor(appConfig: AppConfig, firebaseConfig: FirebaseConfig) {
    this._isDev = appConfig.env === AppEnv.Development;

    // For development
    if (this._isDev) {
      initializeApp(firebaseConfig.client);
    }
  }

  /**
   * Sign in (Development only).
   *
   * @param payload
   * @param response
   */
  @Post('/sign-in')
  async signIn(
    @Body() payload: SignInDto,
    @Res() response: FastifyReply,
  ): Promise<void> {
    if (!this._isDev) {
      response.code(HttpStatus.NOT_FOUND).send();
      return;
    }

    try {
      const credential = await signInWithEmailAndPassword(
        getAuth(),
        payload.email,
        payload.password,
      );

      response.code(HttpStatus.CREATED).send({
        data: credential,
      });
    } catch {
      response.code(HttpStatus.UNAUTHORIZED).send({
        message: 'Incorrect email or password!',
      });
    }
  }
}
