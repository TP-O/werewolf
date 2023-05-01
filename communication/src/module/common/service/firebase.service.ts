import { Injectable } from '@nestjs/common';
import { cert, initializeApp } from 'firebase-admin/app';
import { Auth, getAuth } from 'firebase-admin/auth';
import { FirebaseConfig } from 'src/config';

@Injectable()
export class FirebaseService {
  constructor(config: FirebaseConfig) {
    //Prevent duplicate app initialization
    initializeApp({
      credential: cert({
        projectId: config.productId,
        privateKey: config.privateKey.replace(/\\n/gm, '\n'),
        clientEmail: config.clientEmail,
      }),
      databaseURL: `https://${config.productId}.firebaseio.com`,
    });
  }

  auth(): Auth {
    return getAuth();
  }
}
