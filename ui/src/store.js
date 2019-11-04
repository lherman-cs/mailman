import Vue from "vue";
import Vuex from "vuex";
import { firebase, auth, collections } from "./firebase";
import {
  createSession as _createSession,
  connect as _connect,
  getHost
} from "./api";
import { webSocket } from "rxjs/webSocket";

Vue.use(Vuex);

const store = new Vuex.Store({
  state: {
    user: null,
    session: null,
    mailmanMessage: null,
    currentState: "",
    error: null
  },
  getters: {
    message: state => state.session && state.session.message
  },
  mutations: {
    UPDATE_USER: (state, user) => {
      if (!user) {
        return;
      }

      state.user = {
        email: user.email,
        displayName: user.displayName,
        uid: user.uid,
        photoURL: user.photoURL
      };
    },
    UPDATE_SESSION: (state, session) => (state.session = session),
    UPDATE_MAILMAN_MESSAGE: (state, message) =>
      (state.mailmanMessage = message),
    UPDATE_CURRENT_STATE: (state, currentState) =>
      (state.currentState = currentState),
    UPDATE_ERROR: (state, error) => (state.error = error)
  },
  actions: {
    login,
    logout,
    createSession,
    connect
  }
});

export default store;

const { showError, subscribe, sendMessage } = buildPrivateActions(store);

async function login() {
  try {
    const provider = new firebase.auth.GoogleAuthProvider();
    auth.signInWithRedirect(provider);
    await auth.getRedirectResult();
  } catch (err) {
    showError(err);
  }
}

async function logout() {
  await auth.signOut();
}

async function createSession({ commit }, { password }) {
  try {
    commit("UPDATE_CURRENT_STATE", "Creating a session");
    const session = await _createSession(password);
    commit("UPDATE_SESSION", session);
    commit("UPDATE_CURRENT_STATE", "Waiting for a sender");
    await subscribe();
  } catch (err) {
    showError(err);
  }
}

async function connect({ commit }, { toEmail, password }) {
  try {
    commit("UPDATE_CURRENT_STATE", "Connecting to a session");
    const session = await _connect(toEmail, password);
    commit("UPDATE_SESSION", session);
    await subscribe();
    await sendMessage({ type: "adapter-sender-ready" });
  } catch (err) {
    showError(err);
  }
}

function buildPrivateActions({ commit, state }) {
  return { showError, subscribe, sendMessage };

  function showError(err) {
    commit("UPDATE_ERROR", err);
    commit("UPDATE_CURRENT_STATE", "");
  }

  async function subscribe() {
    const sessionID = state.user.uid;
    const yourSession = collections.sessions.doc(sessionID);
    const host = await getHost();
    const socket = webSocket(`${host}/api/connect`);

    yourSession.onSnapshot(doc => {
      let data = doc.data();
      if (!data) {
        return;
      }

      commit("UPDATE_SESSION", data);
      handleMessage(data.message);
    });

    return;

    async function handleMessage(msg) {
      if (!msg) {
        return;
      }

      switch (msg.type) {
        case "adapter-sender-ready":
          socket.subscribe(handleMailmanMessage);
          await sendMessage({
            type: "adapter-receiver-ready"
          });
          break;
        case "adapter-receiver-ready":
          socket.subscribe(handleMailmanMessage);
          break;
        default:
          socket.next(msg);
          break;
      }
    }

    async function handleMailmanMessage(msg) {
      commit("UPDATE_MAILMAN_MESSAGE", msg);
      if (msg.type === "state") {
        let state = atob(msg.payload);
        // TODO! this is a hacky way to remove double quotes
        state = state.slice(1, state.length - 1);
        commit("UPDATE_CURRENT_STATE", state);
        return;
      }
      await sendMessage({ ...msg });
    }
  }

  async function sendMessage({ type, payload }) {
    if (!state.user) {
      throw new Error("You need to be logged in before you can send messages");
    }

    const session = state.session;
    if (!session) {
      throw new Error("session hasn't been created");
    }

    const to = state.session.from;
    if (!to) {
      throw new Error("the other peer hasn't been found yet");
    }

    const receiverSession = collections.sessions.doc(to.uid);
    await receiverSession.update({
      message: {
        type,
        payload: payload || null
      }
    });
  }
}
