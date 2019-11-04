<template>
  <v-layout justify-center align-center>
    <v-progress-circular
      :rotate="360"
      :size="100"
      :width="15"
      color="primary"
      indeterminate
    ></v-progress-circular>
  </v-layout>
</template>

<script>
export default {
  async mounted() {
    await this.login();
  },
  methods: {
    async login() {
      const nextURL = this.$route.query.nextURL || "/";
      await this.$store.dispatch("login");
      await this.$router.replace(nextURL);
    }
  }
};
</script>
