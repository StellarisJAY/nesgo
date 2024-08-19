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
        <a-button class="bindingBtn" :disabled="createBtnDisable" type="primary" @click="openCreateBindingModal">新建</a-button>
      </a-col>
      <a-col :span="4">
        <a-button class="bindingBtn" :disabled="bindingBtnDisable" type="primary" @click="deleteBinding" danger>删除</a-button>
      </a-col>
    </a-row>

    <a-table :data-source="bindingSelected['bindings']" :columns="bindingColumns" :pagination="false">
      <template #bodyCell="{column, record}">
        <template v-if="column.dataIndex === 'keyboardKey'">
          <KeyboardKeyPicker :limit="1" :buttons="record['buttons']"></KeyboardKeyPicker>
        </template>
        <template v-else-if="column.dataIndex === 'emulatorKey'">
          {{record['emulatorKeyTranslated']}}
        </template>
      </template>
    </a-table>

    <a-row>
      <a-col :span="20"></a-col>
      <a-col :span="4">
        <a-button type="primary" @click="updateBinding">修改</a-button>
      </a-col>
    </a-row>

    <a-modal v-model:open="createBindingModalOpen" title="新建按键绑定">
      <template #footer>
        <a-button>取消</a-button>
        <a-button type="primary" @click="createBinding">创建</a-button>
      </template>
      <p>提示：点击按钮取消绑定，点击‘+’后按下键盘按键添加绑定</p>
      <a-input v-model:value="newBinding.name" placeholder="按键绑定名称"></a-input>
      <a-table :data-source="newBinding.bindings" :columns="bindingColumns" :pagination="false">
        <template #bodyCell="{column, record}">
          <template v-if="column.dataIndex === 'keyboardKey'">
            <KeyboardKeyPicker :limit="1" :buttons="record['buttons']"></KeyboardKeyPicker>
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
import {
  Button,
  Card,
  Col,
  Descriptions,
  Input,
  List,
  message,
  Modal,
  Pagination,
  Row,
  Select,
  Table
} from 'ant-design-vue';
import api from "../api/request.js";
import KeyboardKeyPicker from "./keyboardKeyPicker.vue";


const defaultBindings = [
      {
        "emulatorKey": "Left",
        "emulatorKeyTranslated": "Left",
        "buttons": ["KeyA"],
      },
      {
        "emulatorKey": "Right",
        "emulatorKeyTranslated": "Right",
        "buttons": ["KeyD"],
      },
      {
        "emulatorKey": "Up",
        "emulatorKeyTranslated": "Up",
        "buttons": ["KeyW"],
      },
      {
        "emulatorKey": "Down",
        "emulatorKeyTranslated": "Down",
        "buttons": ["KeyS"],
      },
      {
        "emulatorKey": "A",
        "emulatorKeyTranslated": "A",
        "buttons": ["Space"],
      },
      {
        "emulatorKey": "B",
        "emulatorKeyTranslated": "B",
        "buttons": ["KeyJ"],
      },
      {
        "emulatorKey": "Start",
        "emulatorKeyTranslated": "Start",
        "buttons": ["Enter"],
      },
      {
        "emulatorKey": "Select",
        "emulatorKeyTranslated": "Select",
        "buttons": ["Tab"],
      },
];

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
    KeyboardKeyPicker,
  },
  data() {
    return {
      keyboardBindings: [],      // 用户创建的所有按键绑定列表
      keyboardSelectOptions: [], // 按键绑定选项列表
      bindingSelectedKey: null, // 选中的按键绑定的ID
      bindingSelected: {},      // 选中的按键绑定
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
        bindings: defaultBindings,
      },
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
          const curBindings = _this.keyboardBindings[i]["bindings"];
          for (let j = 0; j < curBindings.length; j++) {
            curBindings[j]["buttons"] = [curBindings[j]["keyboardKey"]];
          }
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
        _this.bindingBtnDisable = false;
      }).catch(_=>{
        message.error("删除失败");
        _this.bindingBtnDisable = false;
      });
    },

    openCreateBindingModal: function() {
      const s = JSON.stringify(defaultBindings);
      this.newBinding = {
        name: "",
        bindings: JSON.parse(s),
      };
      this.createBindingModalOpen = true;
    },

    createBinding: function() {
      const data = this.convertBindingsToApiObj(this.newBinding);
      const _this = this;
      this.createBtnDisable = true;
      this.bindingBtnDisable = true;
      api.post("api/v1/keyboard/binding", data).then(_=>{
        message.success("创建成功");
        _this.listKeyboardBindings();
        _this.createBtnDisable = false;
        _this.bindingBtnDisable = false;
      }).catch(_=>{
        message.error("创建失败");
        _this.createBtnDisable = false;
        _this.bindingBtnDisable = false;
      });
    },

    convertBindingsToApiObj: function(bindingObj) {
      const apiObj = {
        "name": bindingObj.name,
        "bindings": [],
      };
      bindingObj.bindings.forEach(item=>{
        apiObj["bindings"].push({
          "emulatorKey": item.emulatorKey,
          "keyboardKey": item.buttons[0],
        });
      });
      return apiObj;
    },

    updateBinding: function() {
      const data = this.convertBindingsToApiObj(this.bindingSelected);
      data["id"] = this.bindingSelectedKey;
      const _this = this;
      this.createBtnDisable = true;
      this.bindingBtnDisable = true;
      api.put("api/v1/keyboard/binding", data).then(_=>{
        message.success("修改成功");
        _this.listKeyboardBindings();
        _this.createBtnDisable = false;
        _this.bindingBtnDisable = false;
      }).catch(_=>{
        message.error("修改失败");
        _this.createBtnDisable = false;
        _this.bindingBtnDisable = false;
      });
    },
  }
}
</script>

<style>
.bindingBtn {
  width: 90%;
}
</style>