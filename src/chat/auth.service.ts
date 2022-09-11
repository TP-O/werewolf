import { Injectable } from '@nestjs/common';
import { User } from '@prisma/client';
import { Auth } from 'firebase-admin/auth';
import { FirebaseAuth } from 'src/decorator/firebase-auth.decorator';
import { UserId } from 'src/enum/user.enum';
import { PrismaService } from './prisma.service';

@Injectable()
export class AuthService {
  @FirebaseAuth()
  private readonly auth: Auth;

  constructor(private prismaService: PrismaService) {}

  private async getFirebaseId(token: string) {
    let userId: string;

    try {
      const decodedToken = await this.auth.verifyIdToken(token);
      userId = decodedToken.uid;
    } catch {
      userId = '';
    }

    return userId;
  }

  private generateEmptyUser(id: number): User {
    return {
      id,
      fid: '',
      sids: [],
    };
  }

  async getUser(token: string) {
    const fid = await this.getFirebaseId(token);

    if (fid === '') {
      return this.generateEmptyUser(UserId.NonExist);
    }

    const user = await this.prismaService.user.findUnique({
      where: {
        fid,
      },
    });

    return user ?? this.generateEmptyUser(UserId.Asynchronous);
  }
}
