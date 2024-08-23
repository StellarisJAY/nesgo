<template>
  <a-row>
    <a-col :span="12">
      <a-select v-model:value="selectedKey" :options="selections" @change="onSelectionChange"></a-select>
    </a-col>
    <a-col :span="6">
      <a-button>新建</a-button>
    </a-col>
    <a-col :span="6">
      <a-button>删除</a-button>
    </a-col>
  </a-row>
  <a-table :title="_=>'动作列表'" :data-source="selected['actions']" :columns="columns" :pagination="false">
    <template #bodyCell="{column, record}">
      <template v-if="column.dataIndex==='emulatorKey'">
        <a-button>{{record['emulatorKey']}}</a-button>
      </template>
      <template v-else-if="column.dataIndex==='releaseDelay'">
        {{record['releaseDelay']}}ms
      </template>
    </template>
  </a-table>
</template>

<script>
import {Button, message, Table, Select, Row, Col} from 'ant-design-vue';
import api from '../api/request.js';

export default {
  props: {
  },
  components: {
    AButton: Button,
    AButtonGroup: Button.Group,
    ATable: Table,
    ASelect: Select,
    ARow: Row,
    ACol: Col,
  },
  data() {
    return {
      macros: [],
      total: 0,
      selections: [],
      selectedKey: "0",
      selected: {},
      columns: [
        {
          "title": "模拟器按键",
          "dataIndex": "emulatorKey"
        },
        {
          "title": "动作时间",
          "dataIndex": "releaseDelay"
        }
      ]
    }
  },
  created() {
    this.listMacros();
  },
  methods: {
    listMacros: function() {
      const _this = this;
      api.get("api/v1/macros?page=0&pageSize=100").then(resp=>{
        _this.macros = resp['macros'];
        _this.total = resp['total'];
        const macros = resp['macros'];
        const options = [];
        for (let i = 0; i < macros.length; i++) {
          options.push({
            "label": macros[i]['name'],
            "value": macros[i]['id']
          });
        }
        if (options.length > 0) {
          _this.selected = macros[0];
          _this.selectedKey = macros[0]['id'];
          _this.selections = options;
        }
      }).catch(_=>{
        message.error("无法获取用户宏指令列表");
      })
    },
    onSelectionChange: function() {
      this.selected = this.macros.find(item=>item['id'] === this.selectedKey);
    },
  }
}
</script>

<style>
</style>