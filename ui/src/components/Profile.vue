<template>
  <v-card class="mx-auto" max-width="344">
    <v-img
      :src="
        user ? user.photoURL : require('@/assets/img/profile-placeholder.png')
      "
    ></v-img>
    <v-card-title>{{ user ? user.displayName : "????" }}</v-card-title>
    <v-card-subtitle>{{ user ? user.email : "????" }}</v-card-subtitle>
  </v-card>
</template>

<script>
export default {
  props: ["user"]
};
</script>
