<template>
  <div>
    <a-row>
      <a-col :span="10">
        <a-select :disable="bindingBtnDisable" :options="selectOptions" v-model:value="selectedKey"
                  @change="onSelectChange"></a-select>
      </a-col>
      <a-col :span="6"></a-col>
      <a-col :span="4">
        <a-button class="bindingBtn" :disabled="createBtnDisable" :hidden="!allowCreate" type="primary"
                  @click="openCreateBindingModal">新建</a-button>
      </a-col>
      <a-col :span="4">
        <a-button class="bindingBtn" :disabled="deleteBtnDisabled" :hidden="!allowDelete" type="primary" @click="deleteBinding"
                  danger>删除</a-button>
      </a-col>
    </a-row>

    <a-table :data-source="selected['bindings']" :columns="bindingColumns" :pagination="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'keyboardKey'">
          <KeyboardKeyPicker :limit="1" :buttons="record['buttons']" :keycode-translator="translateKeyboardKeyCode"></KeyboardKeyPicker>
        </template>
        <template v-else-if="column.dataIndex === 'emulatorKey'">
          {{ record['emulatorKeyTranslated'] }}
        </template>
      </template>
    </a-table>

    <a-row>
      <a-col :span="20"></a-col>
      <a-col :span="4">
        <a-button type="primary" :disabled="updateBtnDisabled" @click="updateBinding">保存修改</a-button>
      </a-col>
    </a-row>

    <a-modal v-model:open="createBindingModalOpen" title="新建按键绑定">
      <template #footer>
        <a-button type="primary" @click="createBinding">创建</a-button>
      </template>
      <p>提示：点击按钮取消绑定，点击‘+’后按下键盘按键添加绑定</p>
      <a-input v-model:value="newBinding.name" placeholder="按键绑定名称"></a-input>
      <a-table :data-source="newBinding.bindings" :columns="bindingColumns" :pagination="false">
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'keyboardKey'">
            <KeyboardKeyPicker :limit="1" :buttons="record['buttons']" :keycode-translator="translateKeyboardKeyCode"></KeyboardKeyPicker>
          </template>
          <template v-else-if="column.dataIndex === 'emulatorKey'">
            {{ record['emulatorKeyTranslated'] }}
          </template>
        </template>
      </a-table>
    </a-modal>
  </div>
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

const defaultSettings = {
  "id": "0",
  "name": "默认设置",
  "bindings": [
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
  ]
}
const keyboardKeyTranslations = {
  "KeyA": "A",
  "KeyB": "B",
  "KeyC": "C",
  "KeyD": "D",
  "KeyE": "E",
  "KeyF": "F",
  "KeyG": "G",
  "KeyH": "H",
  "KeyI": "I",
  "KeyJ": "J",
  "KeyK": "K",
  "KeyL": "L",
  "KeyM": "M",
  "KeyN": "N",
  "KeyO": "O",
  "KeyP": "P",
  "KeyQ": "Q",
  "KeyR": "R",
  "KeyS": "S",
  "KeyT": "T",
  "KeyU": "U",
  "KeyV": "V",
  "KeyW": "W",
  "KeyX": "X",
  "KeyY": "Y",
  "KeyZ": "Z",
  "Space": "空格",
  "ArrowLeft": "←",
  "ArrowRight": "→",
  "ArrowUp": "↑",
  "ArrowDown": "↓",
}
export default {
  props: {
    showDefault: Boolean,
    allowCreate: Boolean,
    allowDelete: Boolean,
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
      keyboardSettings: [],      // 用户创建的所有按键绑定列表
      selectOptions: [], // 按键绑定选项列表
      selectedKey: null, // 选中的按键绑定的ID
      selected: {},      // 选中的按键绑定
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
      bindingBtnDisable: false,
      deleteBtnDisabled: false,
      updateBtnDisabled: false,
      createBtnDisable: false,
      createBindingModalOpen: false,
      newBinding: {},
    }
  },
  created() {
    this.listKeyboardBindings();
  },
  methods: {
    listKeyboardBindings: function () {
      const _this = this;
      const showDefault = this.showDefault;
      api.get("api/v1/keyboard/bindings?page=0&pageSize=100").then(resp => {
        _this.keyboardSettings = resp["bindings"];
        _this.total = resp["total"];
        let options = [];
        for (let i = 0; i < _this.keyboardSettings.length; i++) {
          options.push({
            value: _this.keyboardSettings[i]["id"],
            label: _this.keyboardSettings[i]["name"],
          });
          const curSetting = _this.keyboardSettings[i]["bindings"];
          // 每个配置的按键列表（为了将来支持多按建绑定）
          for (let j = 0; j < curSetting.length; j++) {
            curSetting[j]["buttons"] = [curSetting[j]["keyboardKey"]];
          }
        }
        // 添加默认配置
        if(showDefault) {
          _this.keyboardSettings.push(JSON.parse(JSON.stringify(defaultSettings)));
          options.push({value: defaultSettings["id"], label: defaultSettings["name"]});
        }
        if(options.length === 0) {
          _this.selected = {};
          _this.selectedKey = "0";
          _this.selectOptions = [];
          return;
        }
        _this.selected = _this.keyboardSettings[0];
        _this.selectedKey = _this.keyboardSettings[0]["id"];
        _this.selectOptions = options;
        // 禁止在列表删除或修改默认配置
        if (_this.selectedKey === "0") {
          _this.deleteBtnDisabled = true;
          _this.updateBtnDisabled = true;
        }
      }).catch(_ => {
        message.error("获取按键绑定失败");
      });
    },

    onSelectChange: function (ev) {
      this.selected = this.keyboardSettings.find(item => item["id"] === this.selectedKey);
      this.deleteBtnDisabled = this.selectedKey === '0';
      this.updateBtnDisabled = this.selectedKey === '0';
    },

    deleteBinding: function () {
      const _this = this;
      this.bindingBtnDisable = true;
      api.delete("api/v1/keyboard/binding/" + this.selectedKey).then(_ => {
        message.success("删除成功");
        _this.listKeyboardBindings();
        _this.bindingBtnDisable = false;
      }).catch(_ => {
        message.error("删除失败");
        _this.bindingBtnDisable = false;
      });
    },

    openCreateBindingModal: function () {
      const s = JSON.stringify(defaultSettings);
      this.newBinding = JSON.parse(s);
      this.createBindingModalOpen = true;
    },

    createBinding: function () {
      if (!this.verifySetting(this.newBinding)) return;
      const data = this.convertBindingsToApiObj(this.newBinding);
      const _this = this;
      this.createBtnDisable = true;
      this.bindingBtnDisable = true;
      api.post("api/v1/keyboard/binding", data).then(_ => {
        message.success("创建成功");
        _this.listKeyboardBindings();
        _this.createBtnDisable = false;
        _this.bindingBtnDisable = false;
      }).catch(_ => {
        message.error("创建失败");
        _this.createBtnDisable = false;
        _this.bindingBtnDisable = false;
      });
    },

    convertBindingsToApiObj: function (bindingObj) {
      const apiObj = {
        "name": bindingObj.name,
        "bindings": [],
      };
      bindingObj.bindings.forEach(item => {
        apiObj["bindings"].push({
          "emulatorKey": item.emulatorKey,
          "keyboardKey": item.buttons[0],
        });
      });
      return apiObj;
    },

    updateBinding: function () {
      if (!this.verifySetting(this.selected)) return;
      const data = this.convertBindingsToApiObj(this.selected);
      data["id"] = this.selectedKey;
      const _this = this;
      this.createBtnDisable = true;
      this.bindingBtnDisable = true;
      api.put("api/v1/keyboard/binding", data).then(_ => {
        message.success("修改成功");
        _this.listKeyboardBindings();
        _this.createBtnDisable = false;
        _this.bindingBtnDisable = false;
      }).catch(_ => {
        message.error("修改失败");
        _this.createBtnDisable = false;
        _this.bindingBtnDisable = false;
      });
    },

    translateKeyboardKeyCode: function(code) {
      if (keyboardKeyTranslations[code]) {
        return keyboardKeyTranslations[code];
      } else {
        return code;
      }
    },

    verifySetting: function(setting) {
      let dict = {};
      if (setting.name === '') {
        message.error("按键绑定名称不能为空");
        return false;
      }
      if (setting.name === defaultSettings.name) {
        message.error("按键绑定名称不能为“默认绑定”");
        return false;
      }
      for(let i = 0; i < setting.bindings.length; i++) {
        let item = setting.bindings[i];
        if (item.buttons.length === 0) {
          message.error("模拟器按键" + item.emulatorKey + "未设置")
          return false;
        }
        if (dict[item.buttons[0]]) {
          message.error("键盘按键冲突！");
          return false;
        }
        dict[item.buttons[0]] = true;
      }
      return true;
    },
  }
}
</script>

<style>
.bindingBtn {
  width: 90%;
}
</style>