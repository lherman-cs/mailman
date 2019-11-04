<template>
  <v-layout justify-center align-center>
    <my-progress
      v-show="currentState"
      :sender="user"
      :receiver="session ? session.from : null"
      :message="currentState"
    ></my-progress>
    <v-form @submit.prevent="connect()" v-show="!currentState">
      <v-container>
        <v-row class="align-center">
          <v-col cols="4">
            <v-text-field
              label="receiver email"
              v-model="toEmail"
              required
            ></v-text-field>
          </v-col>
          <v-col cols="4">
            <v-text-field
              label="session password"
              v-model="password"
              type="password"
              required
            ></v-text-field>
          </v-col>
          <v-col cols="4">
            <v-btn color="primary" type="submit">connect</v-btn>
          </v-col>
        </v-row>
      </v-container>
    </v-form>
  </v-layout>
</template>

<script>
import Progress from "@/components/Progress.vue";
import { mapState } from "vuex";

export default {
  components: {
    MyProgress: Progress
  },
  data() {
    return {
      toEmail: "",
      password: ""
    };
  },
  computed: {
    ...mapState(["user", "session", "currentState"])
  },
  methods: {
    async connect() {
      await this.$store.dispatch("connect", {
        toEmail: this.toEmail,
        password: this.password
      });
    }
  }
};
</script>
