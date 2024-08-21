<template>
  <a-card :bordered="false">
    <KeyboardSetting :allow-delete="true" :allow-create="true" :show-default="false"></KeyboardSetting>
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
    }
  },
}
</script>

<style>
</style>