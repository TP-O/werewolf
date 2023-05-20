import { initializeApp } from 'firebase/app'
import { GoogleAuthProvider, signOut as fSignOut, getAuth, signInWithEmailAndPassword, signInWithPopup } from 'firebase/auth'

const firebaseConfig = {
  apiKey: 'AIzaSyBDEGHwraskC2U96zUEy5HwN3LFBlZNPDE',
  authDomain: 'werewolf-fdac9.firebaseapp.com',
  projectId: 'werewolf-fdac9',
  storageBucket: 'werewolf-fdac9.appspot.com',
  messagingSenderId: '244960264403',
  appId: '1:244960264403:web:a2a0fe9dd6b00495481426',
  measurementId: 'G-94F5EE6YCY',
}

const app = initializeApp(firebaseConfig)
const firebaseAuth = getAuth(app)
firebaseAuth.useDeviceLanguage()

const provider = new GoogleAuthProvider()
provider.addScope('https://www.googleapis.com/auth/contacts.readonly')

export async function signIn(email: string, password: string) {
  try {
    await signInWithEmailAndPassword(firebaseAuth, email, password)
  }
  catch (err: any) {
    if (err.code === 'auth/user-not-found')
      throw new Error('Email or password is incorrect.')
    else
      throw new Error('Unkown error.')
  }
}

export async function signInWithGoogle() {
  try {
    await signInWithPopup(firebaseAuth, provider)
  }
  catch (err) {
    console.error(err)
  }
}

export async function signOut() {
  await fSignOut(firebaseAuth)
}

export { firebaseAuth }
