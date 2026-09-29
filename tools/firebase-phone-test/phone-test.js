import { initializeApp } from "firebase/app";
import {
  getAuth,
  RecaptchaVerifier,
  signInWithPhoneNumber,
} from "firebase/auth";

const firebaseConfig = {
  apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
  authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
  projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID,
  appId: import.meta.env.VITE_FIREBASE_APP_ID,
};

const phoneNumber = import.meta.env.VITE_TEST_PHONE_NUMBER;
const verificationCode = import.meta.env.VITE_TEST_VERIFICATION_CODE;

const app = initializeApp(firebaseConfig);
const auth = getAuth(app);

// Testing only - allows Firebase fictional phone numbers.
auth.settings.appVerificationDisabledForTesting = true;

async function loginWithTestPhone() {
  try {
    const recaptchaVerifier = new RecaptchaVerifier(
      auth,
      "recaptcha-container",
      {
        size: "invisible",
      }
    );

    const confirmationResult = await signInWithPhoneNumber(
      auth,
      phoneNumber,
      recaptchaVerifier
    );

    const userCredential = await confirmationResult.confirm(
      verificationCode
    );

    const user = userCredential.user;

    console.log("✅ Firebase UID:");
    console.log(user.uid);

    const idToken = await user.getIdToken();

    console.log("✅ Firebase ID Token:");
    console.log(idToken);
  } catch (error) {
    console.error("❌ Firebase Phone Auth failed:");
    console.error(error);
  }
}

loginWithTestPhone();