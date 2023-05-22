import { initializeApp } from 'firebase/app'
import {
  GoogleAuthProvider,
  createUserWithEmailAndPassword,
  signOut as fSignOut, getAuth,
  signInWithEmailAndPassword,
  signInWithPopup,
} from 'firebase/auth'
import { error } from 'loglevel'

const firebaseConfig = JSON.parse(import.meta.env.VITE_FIREBASE_CONFIG)
const app = initializeApp(firebaseConfig)
const firebaseAuth = getAuth(app)
firebaseAuth.useDeviceLanguage()

const ggProvider = new GoogleAuthProvider()
ggProvider.addScope('https://www.googleapis.com/auth/contacts.readonly')

let isAuthStateChanged = false

/**
 * Use to wait for firebase authetication at the time
 * the app is launched.
 */
function waitForAuthState(): Promise<void> {
  return new Promise<void>((resolve) => {
    if (isAuthStateChanged)
      return resolve()

    isAuthStateChanged = true
    firebaseAuth.onAuthStateChanged(() => resolve())
  })
}

async function getIdToken(): Promise<string | undefined> {
  await waitForAuthState()
  return firebaseAuth.currentUser?.getIdToken()
}

async function signUp(email: string, password: string): Promise<void> {
  try {
    await createUserWithEmailAndPassword(firebaseAuth, email, password)
  }
  catch (err: any) {
    switch (err.code) {
      case 'auth/email-already-in-use':
        throw new Error('Email is already in use')

      default:
        error('Sign-up error:', err.message)
        throw new Error('Please try again')
    }
  }
}

async function signIn(email: string, password: string): Promise<void> {
  try {
    await signInWithEmailAndPassword(firebaseAuth, email, password)
  }
  catch (err: any) {
    switch (err.code) {
      case 'auth/user-not-found':
      case 'auth/wrong-password':
        throw new Error('Email or password is incorrect')

      default:
        error('Sign-in error:', err.message)
        throw new Error('Something went wrong')
    }
  }
}

async function signInWithGoogle(): Promise<void> {
  try {
    await signInWithPopup(firebaseAuth, ggProvider)
  }
  catch (err: any) {
    error('Sign-in error:', err.message)
    throw new Error('Please try another way to sign in')
  }
}

async function signOut(): Promise<void> {
  try {
    await waitForAuthState()
    await fSignOut(firebaseAuth)
  }
  catch (err: any) {
    error('Sign-out error:', err.message)
    throw new Error('Unable to sign out')
  }
}

export const auth = {
  waitForAuthState,
  getIdToken,
  signUp,
  signIn,
  signInWithGoogle,
  signOut,
  raw: firebaseAuth,
}
