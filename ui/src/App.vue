<template>
  <v-app id="”app”">
    <v-content>
      <v-container fill-height>
        <router-view />
      </v-container>
    </v-content>
    <v-snackbar v-model="showError" color="error">
      <b>{{ error && error.message }}</b>
      <v-btn text @click="showError = false">Close</v-btn>
    </v-snackbar>
  </v-app>
</template>

<script>
import { mapState } from "vuex";

export default {
  data() {
    return {
      showError: false
    };
  },
  computed: {
    ...mapState(["error"])
  },
  watch: {
    error() {
      console.log("ERROR");
      this.showError = true;
    }
  }
};
</script>
