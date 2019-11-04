import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import { auth } from "./firebase";
import vuetify from "@/plugins/vuetify";

Vue.config.productionTip = false;

let vue = null;
auth.onAuthStateChanged(async function(user) {
  console.log(user);
  await store.commit("UPDATE_USER", user);

  if (vue) {
    return;
  }

  vue = new Vue({
    router,
    store,
    vuetify,
    render: h => h(App)
  }).$mount("#app");
});
