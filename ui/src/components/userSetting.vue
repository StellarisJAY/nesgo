<template>
  <a-card :bordered="false">
    <a-row>
      <a-col :span="10">
        <a-select :disable="bindingBtnDisable" :options="keyboardSelectOptions" v-model:value="bindingSelectedKey" @change="onKeyboardBindingSelectChange"></a-select>
      </a-col>
      <a-col :span="6">
        <a-button class="bindingBtn" :disabled="bindingBtnDisable" type="primary">恢复默认</a-button>
      </a-col>
      <a-col :span="4">
        <a-button class="bindingBtn" :disabled="createBtnDisable" type="primary" @click="_=>{createBindingModalOpen=true;}">新建</a-button>
      </a-col>
      <a-col :span="4">
        <a-button class="bindingBtn" :disabled="bindingBtnDisable" type="primary" @click="deleteBinding" danger>删除</a-button>
      </a-col>
    </a-row>

    <a-table :data-source="bindingSelected['bindings']" :columns="bindingColumns" :pagination="false">
      <template #bodyCell="{column, record}">
        <template v-if="column.dataIndex === 'keyboardKey'">
          {{record['keyboardKeyTranslated']}}
        </template>
        <template v-else-if="column.dataIndex === 'emulatorKey'">
          {{record['emulatorKeyTranslated']}}
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="createBindingModalOpen" title="新建按键绑定">
      <template #footer>
        <a-button>取消</a-button>
        <a-button type="primary">创建</a-button>
      </template>
      <a-input v-model:value="newBinding.name" placeholder="按键绑定名称"></a-input>
      <a-table :data-source="defaultBindings" :columns="bindingColumns" :pagination="false">
        <template #bodyCell="{column, record}">
          <template v-if="column.dataIndex === 'keyboardKey'">
            {{record['keyboardKeyTranslated']}}
          </template>
          <template v-else-if="column.dataIndex === 'emulatorKey'">
            {{record['emulatorKeyTranslated']}}
          </template>
        </template>
      </a-table>
    </a-modal>
  </a-card>
</template>

<script>
import { List, Button, Descriptions, Pagination, Table, Card, Select, Modal, Input} from 'ant-design-vue';
import { Row, Col } from "ant-design-vue";
import {message} from "ant-design-vue";
import api from "../api/request.js";

export default {
  props: {
    roomId: String,
  },
  components: {
    AButton: Button,
    ARow: Row,
    ACol: Col,
    AList: List,
    AListItem: List.Item,
    ADescriptions: Descriptions,
    ADescriptionsItem: Descriptions.Item,
    APagination: Pagination,
    ATable: Table,
    ACard: Card,
    ASelect: Select,
    AModal: Modal,
    AInput: Input,
  },
  data() {
    return {
      keyboardBindings: [],
      keyboardSelectOptions: [],
      bindingSelectedKey: null,
      bindingSelected: {},
      bindingColumns: [
        {
          "title": "键盘按键",
          "dataIndex": "keyboardKey"
        },
        {
          "title": "模拟器按键",
          "dataIndex": "emulatorKey"
        }
      ],
      page: 1,
      pageSize: 100,
      total: 0,
      bindingBtnDisable: false,
      createBtnDisable: false,
      createBindingModalOpen: false,
      newBinding: {
        name: "",
        bindings: this.defaultBindings,
      },

      defaultBindings: [
        {
          "keyboardKey": "KeyA",
          "emulatorKey": "Left",
          "keyboardKeyTranslated": "A",
          "emulatorKeyTranslated": "Left"
        },
        {
          "keyboardKey": "KeyD",
          "emulatorKey": "Right",
          "keyboardKeyTranslated": "D",
          "emulatorKeyTranslated": "Right"
        },
        {
          "keyboardKey": "KeyW",
          "emulatorKey": "Up",
          "keyboardKeyTranslated": "W",
          "emulatorKeyTranslated": "Up"
        },
        {
          "keyboardKey": "KeyS",
          "emulatorKey": "Down",
          "keyboardKeyTranslated": "S",
          "emulatorKeyTranslated": "Down"
        },
        {
          "keyboardKey": "KeyJ",
          "emulatorKey": "B",
          "keyboardKeyTranslated": "J",
          "emulatorKeyTranslated": "B"
        },
        {
          "keyboardKey": "Space",
          "emulatorKey": "A",
          "keyboardKeyTranslated": "Space",
          "emulatorKeyTranslated": "A"
        },
        {
          "keyboardKey": "Enter",
          "emulatorKey": "Start",
          "keyboardKeyTranslated": "Enter",
          "emulatorKeyTranslated": "Start"
        },
        {
          "keyboardKey": "Tab",
          "emulatorKey": "Select",
          "keyboardKeyTranslated": "Tab",
          "emulatorKeyTranslated": "Select"
        }
      ],

    }
  },
  created() {
    this.listKeyboardBindings();
  },
  methods: {
    listKeyboardBindings: function() {
      const _this = this;
      api.get("api/v1/keyboard/bindings?page=0&pageSize=100").then(resp=>{
        _this.keyboardBindings = resp["bindings"];
        _this.total = resp["total"];
        let options = [];
        for (let i = 0; i < _this.keyboardBindings.length; i++) {
          options.push({
            value: _this.keyboardBindings[i]["id"],
            label: _this.keyboardBindings[i]["name"],
          });
        }
        if (options.length === 0) {
          _this.bindingBtnDisable = true;
          _this.keyboardSelectOptions = [];
          _this.bindingSelectedKey = "";
          _this.bindingSelected = {};
        }else {
          _this.bindingSelectedKey = options[0]["value"];
          _this.bindingSelected = _this.keyboardBindings[0];
          _this.keyboardSelectOptions = options;
        }
      }).catch(_=>{
        message.error("获取按键绑定失败");
      });
    },

    onKeyboardBindingSelectChange: function(ev) {
      this.bindingSelected = this.keyboardBindings.find(item=>item["id"]===this.bindingSelectedKey);
    },

    deleteBinding: function() {
      const _this = this;
      this.bindingBtnDisable = true;
      api.delete("api/v1/keyboard/binding/"+this.bindingSelectedKey).then(_=>{
        message.success("删除成功");
        _this.listKeyboardBindings();
      }).catch(_=>{
        message.error("删除失败");
        _this.bindingBtnDisable = false;
      });
    },

    createBinding: function() {
      const _this = this;
      this.createBtnDisable = true;
      this.bindingBtnDisable = true;
      api.post("api/v1/keyboard/binding", this.newBinding).then(_=>{
        message.success("创建成功");
        _this.listKeyboardBindings();
        _this.createBtnDisable = false;
        _this.bindingBtnDisable = false;
      }).catch(_=>{
        message.error("创建失败");
        _this.createBtnDisable = false;
        _this.bindingBtnDisable = false;
      })
    },
  }
}
</script>

<style>
.bindingBtn {
  width: 90%;
}
</style>