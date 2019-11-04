<template>
  <v-layout justify-center align-center>
    <my-progress
      v-show="currentState"
      :sender="session ? session.from : null"
      :receiver="user"
      :message="currentState"
    ></my-progress>
    <v-form @submit.prevent="createSession()" v-show="!currentState">
      <v-container>
        <v-row class="align-center">
          <v-col cols="6">
            <v-text-field
              label="session password"
              type="password"
              v-model="password"
              required
            ></v-text-field>
          </v-col>
          <v-col cols="6">
            <v-btn color="primary" type="submit">create session</v-btn>
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
      password: ""
    };
  },
  computed: {
    ...mapState(["user", "session", "currentState"])
  },
  async mounted() {},
  methods: {
    async createSession() {
      await this.$store.dispatch("createSession", { password: this.password });
    }
  }
};
</script>
