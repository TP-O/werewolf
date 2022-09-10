import { Injectable } from '@nestjs/common';
import { cert, getApp, getApps, initializeApp } from 'firebase-admin/app';
import { Auth, getAuth } from 'firebase-admin/auth';
import { env } from 'process';

@Injectable()
export class AuthService {
  private app: Auth;

  constructor() {
    if (getApps().length === 0) {
      initializeApp({
        credential: cert({
          projectId: env.FIREBASE_PRODUCT_ID,
          privateKey: String(process.env.FIREBASE_PRIVATE_KEY).replace(
            /\\n/g,
            '\n',
          ),
          clientEmail: env.FIREBASE_CLIENT_EMAIL,
        }),
        databaseURL: `https://${env.FIREBASE_PRODUCT_ID}.firebaseio.com`,
      });
    }

    this.app = getAuth(getApp());
  }

  async getUserId(token: string) {
    let userId: string;

    try {
      const decodedToken = await this.app.verifyIdToken(token);
      userId = decodedToken.uid;
    } catch {
      userId = '';
    }

    return userId;
  }
}
