<template>
  <a-button-group>
    <a-button v-for="b in buttons" @click="deleteButton(b)">{{b}}</a-button>
    <a-button @click="addButton" :hidden="limit === buttons.length">+</a-button>
  </a-button-group>
</template>

<script>
import {Button, message} from 'ant-design-vue';

export default {
  props: {
    limit: Number,
    buttons: Array,
  },
  components: {
    AButton: Button,
    AButtonGroup: Button.Group,
  },
  data() {
    return {
    }
  },
  created() {
  },
  methods: {
    addButton: function() {
      if (this.buttons.length === this.limit) {
        message.warn("最多绑定"+this.limit+"个按键");
        return;
      }
      addEventListener("keyup", this.keyUpListener, false);
    },

    keyUpListener: function(ev) {
      removeEventListener("keyup", this.keyUpListener, false);
      const idx = this.buttons.findIndex(item=>item===ev.code);
      if (idx === -1 && this.buttons.length < this.limit) {
        this.buttons.push(ev.code);
      }
    },

    deleteButton: function(b) {
      const idx = this.buttons.findIndex(item=>item===b);
      this.buttons.splice(idx, 1);
    },
  }
}
</script>

<style>
</style>