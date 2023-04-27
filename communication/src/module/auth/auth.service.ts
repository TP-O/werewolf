import {
  Injectable,
  InternalServerErrorException,
  UnauthorizedException,
} from '@nestjs/common';
import { Auth } from 'firebase-admin/auth';
import { FirebaseService, PrismaService } from '../common';

@Injectable()
export class AuthService {
  private readonly _auth: Auth;

  constructor(
    private prismaService: PrismaService,
    firebaseService: FirebaseService,
  ) {
    this._auth = firebaseService.auth();
  }

  /**
   * Get uid generated by firebase authentication using
   * ID token provided by it.
   *
   * @param token
   * @returns
   */
  private async getFirebaseUserId(token: string) {
    let fid: string;

    try {
      const decodedToken = await this._auth.verifyIdToken(token);
      fid = decodedToken.uid;
    } catch {
      throw new UnauthorizedException('Invalid access token!');
    }

    return fid;
  }

  /**
   * Get a corresponding player on the entered token. Throw an error
   * if authentication failed.
   *
   * @param token ID token provided by firebase authentication.
   * @returns
   */
  async getPlayer(token: string) {
    const id = await this.getFirebaseUserId(token);
    const player = await this.prismaService.player.findUnique({
      where: {
        id,
      },
    });
    if (!player) {
      throw new InternalServerErrorException(
        'Please connect again after a while!',
      );
    }

    return player;
  }
}
