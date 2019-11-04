<template>
  <v-container>
    <v-row class="align-center">
      <v-col cols="4">
        <profile :user="sender"></profile>
      </v-col>
      <v-col cols="4">
        <div class="text-center">
          <h1>{{ message }}</h1>
          <v-progress-circular
            indeterminate
            color="primary"
            size="48"
            width="6"
            class="ma-5"
          ></v-progress-circular>
        </div>
      </v-col>
      <v-col cols="4">
        <profile :user="receiver"></profile>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import Profile from "./Profile.vue";

export default {
  components: {
    Profile
  },
  props: ["sender", "receiver", "message"]
};
</script>
