import { functions } from "./firebase";

const _createSession = functions.httpsCallable("createSession");
const _connect = functions.httpsCallable("connect");

export function createSession(password) {
  return _createSession({ password }).then(result => result.data);
}

export function connect(toEmail, password) {
  return _connect({ toEmail, password }).then(result => result.data);
}

export function getHost() {
  return fetch("/api/getHost").then(res => res.text());
}

export function getPeerType() {
  return fetch("/api/getPeerType").then(res => res.text());
}

export const mailman = {
  close: () => null
};
