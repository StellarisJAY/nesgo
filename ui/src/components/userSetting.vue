<template>
  <a-card :bordered="false">
    <KeyboardSetting v-if="!mobileDevice" :allow-delete="true" :allow-create="true" :show-default="false"></KeyboardSetting>
  </a-card>
</template>

<script>
import {
  Button,
  Card,
  Col,
  Row,
} from 'ant-design-vue';
import KeyboardSetting from "./keyboardSetting.vue";
import platform from "../util/platform.js";
export default {
  props: {
  },
  components: {
    AButton: Button,
    ARow: Row,
    ACol: Col,
    ACard: Card,
    KeyboardSetting
  },
  data() {
    return {
      mobileDevice: false
    }
  },
  created() {
    this.mobileDevice = platform.isMobile();
  }
}
</script>

<style>
</style>