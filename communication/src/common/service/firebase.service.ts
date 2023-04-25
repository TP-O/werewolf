import { Injectable } from '@nestjs/common';
import { cert, getApps, initializeApp } from 'firebase-admin/app';
import { Auth, getAuth } from 'firebase-admin/auth';
import { FirebaseConfig } from 'src/config/firebase';

@Injectable()
export class FirebaseService {
  constructor(config: FirebaseConfig) {
    if (getApps().length === 0) {
      initializeApp({
        credential: cert({
          projectId: config.productId,
          privateKey: config.privateKey.replace(/\\n/gm, '\n'),
          clientEmail: config.clientEmail,
        }),
        databaseURL: `https://${config.productId}.firebaseio.com`,
      });
    }
  }

  auth(): Auth {
    return getAuth();
  }
}
