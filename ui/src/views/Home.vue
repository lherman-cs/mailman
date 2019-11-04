<template>
  <sender v-if="peerType === 'sender'"></sender>
  <receiver v-else-if="peerType === 'receiver'"></receiver>
</template>

<script>
import { getPeerType } from "../api";
import Sender from "./Sender.vue";
import Receiver from "./Receiver.vue";

export default {
  components: {
    Sender,
    Receiver
  },
  data() {
    return {
      peerType: ""
    };
  },
  async mounted() {
    this.peerType = await getPeerType();
  }
};
</script>
