import * as firebase from "firebase/app";
import "firebase/auth";
import "firebase/functions";
import "firebase/firestore";

// Your web app's Firebase configuration
const firebaseConfig = {
  apiKey: "AIzaSyAAp-IjJ33dlTdfFry8ikpl2mKbxMjKWnk",
  authDomain: "mailman-33702.firebaseapp.com",
  databaseURL: "https://mailman-33702.firebaseio.com",
  projectId: "mailman-33702",
  storageBucket: "",
  messagingSenderId: "115301841410",
  appId: "1:115301841410:web:c61092f61cb8ae36a1851b"
};
// Initialize Firebase
firebase.initializeApp(firebaseConfig);

export { default as firebase } from "firebase/app";
export const auth = firebase.auth();
export const functions = firebase.functions();
export const firestore = firebase.firestore();
export const collections = {
  sessions: firestore.collection("sessions")
};
