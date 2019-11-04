import * as functions from 'firebase-functions';
import * as admin from 'firebase-admin';
import { SHA256 } from "crypto-js";
import { Session, User } from './model';

// The Firebase Admin SDK to access the Firebase Realtime Database.
admin.initializeApp(functions.config().firebase);
const db = admin.firestore();
const sessions = db.collection("sessions");
const HttpsError = functions.https.HttpsError;

export const createSession = functions.https.onCall(async (data, context) => {
    const [password]: string[] = required(data, ["password"]);
    const owner = contextToUser(context);

    const hashedPassword = SHA256(password).toString();
    const doc = sessions.doc(owner.uid);
    const session: Session = {
        owner,
        password: hashedPassword,
        from: null,
        message: null
    }
    try {
        await doc.set(session);
    } catch (e) {
        throw new HttpsError("internal", e.message);
    }

    return session;
});

export const connect = functions.https.onCall(async (data, context) => {
    const [toEmail, password]: string[] = required(data, ["toEmail", "password"]);
    const from = contextToUser(context);

    const query = sessions.where("owner.email", "==", toEmail).limit(1);
    let results: FirebaseFirestore.QuerySnapshot;
    try {
        results = await query.get();
    } catch (e) {
        throw new HttpsError("internal", e.message);
    }

    if (results.empty) {
        throw new HttpsError("out-of-range", `${toEmail} either hasn't registered or run the program`);
    }

    const result = results.docs[0];
    const resultData = result.data() as Session;
    if (resultData.from) {
        throw new HttpsError("out-of-range", `${toEmail} hasn't run the program`);
    }

    const resultHashedPassword = resultData.password;
    const hashedPassword = SHA256(password).toString();
    if (hashedPassword !== resultHashedPassword) {
        throw new HttpsError("permission-denied", "your password is incorrect");
    }

    const resultOwner = resultData.owner;
    const doc = sessions.doc(from.uid);
    const session: Session = {
        owner: from,
        from: resultOwner,
        password: hashedPassword,
        message: null
    }
    try {
        await doc.set(session);
    } catch (e) {
        throw new HttpsError("internal", e.message);
    }

    try {
        await result.ref.update({ from });
    } catch (e) {
        throw new HttpsError("internal", e.message);
    }

    return session;
});

function required(data: any, requiredArgs: string[]): any[] {
    const args = [];
    const missing = [];
    for (const requiredArg of requiredArgs) {
        const value = data[requiredArg];
        if (!value) {
            missing.push(requiredArg);
            continue;
        }
        args.push(value);
    }

    if (missing.length > 0) {
        throw new HttpsError("invalid-argument", `required: ${missing.join(",")}`);
    }

    return args;
}

function contextToUser(context: functions.https.CallableContext): User {
    const token = context.auth!.token;

    return {
        uid: token.uid,
        email: token.email,
        displayName: token.name,
        photoURL: token.picture
    };
}