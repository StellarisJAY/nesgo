<template>
  <a-card :bordered="false">
    <!--服务器列表-->
    <a-list :grid="{ gutter: 16, xs: 1, sm: 2, md: 2, lg: 2, xl: 3, xxl: 3 }" :data-source="endpoints">
      <template #renderItem="{item}">
        <a-list-item>
          <a-card :title="item['id']">
            <template #actions>
              <a-button type="link" @click="_=>openEndpointDetail(item)">详情</a-button>
            </template>
            <ul style="text-align: left">
              <li>地址: {{item["address"]}}</li>
              <li>CPU用量：{{item["cpuUsage"]}}%</li>
              <li>内存使用：{{ formatMemory(item["memoryUsed"]) }}/{{formatMemory(item["memoryTotal"])}}</li>
              <li>模拟器数量：{{item["emulatorCount"]}}</li>
            </ul>
          </a-card>
        </a-list-item>
      </template>
    </a-list>
    <!--分页-->
    <a-pagination v-model:current="page" :total="total" v-model:pageSize="pageSize" @change="onPageChange" />

    <!--活跃房间列表-->
    <a-modal :open="endpointDetailOpen" title="活跃房间列表" @cancel="_=>{endpointDetailOpen=false;}" width="100%">
      <template #footer></template>
      <a-table :data-source="activeRooms" :columns="activeRoomsTableColumns" :pagination="false">
        <template #bodyCell="{column, record}">
          <template v-if="column.dataIndex==='members'">
            {{record["memberCount"]}}/{{record["memberLimit"]}}
          </template>
          <template v-else-if="column.dataIndex==='connections'">
            {{record["activeConnections"]}}/{{record["connections"]}}
          </template>
          <template v-else-if="column.dataIndex==='uptime'">
            {{formatUptime(record["uptime"])}}
          </template>
          <template v-else-if="column.dataIndex==='operations'">
            <a-popconfirm title="确认关闭该房间？" ok-text="确认" cancel-text="取消" @confirm="_=>closeRoom(record['roomId'])">
              <a-button type="link" danger>关闭</a-button>
            </a-popconfirm>
            <a-popconfirm title="确认删除该房间？" ok-text="确认" cancel-text="取消" @confirm="_=>deleteRoom(record['roomId'])">
              <a-button type="link" danger>删除</a-button>
            </a-popconfirm>
          </template>
        </template>
      </a-table>
    </a-modal>

  </a-card>
</template>

<script>
import { Card, Button, List, Modal, Pagination, Table, Popconfirm} from 'ant-design-vue';
import { Row, Col } from "ant-design-vue";
import {message} from "ant-design-vue";
import api from "../api/request.js";
import {InboxOutlined} from "@ant-design/icons-vue";

export default {
  props: {
  },
  components: {
    InboxOutlined,
    ARow: Row,
    ACol: Col,
    ACard: Card,
    AButton: Button,
    AList: List,
    AListItem: List.Item,
    AModal:  Modal,
    APagination: Pagination,
    ATable: Table,
    APopconfirm: Popconfirm,
  },
  data() {
    return {
      endpoints: [
        {"id": "abc", "address": "localhost:4030"}
      ],
      page: 1,
      pageSize: 10,
      total: 0,

      endpointDetailOpen: false,
      selectedEndpoint: {},
      activeRooms: [],

      activeRoomsTableColumns: [
        {
          "title": "房间名",
          "dataIndex": "name",
          "key": "name",
        },
        {
          "title": "房主",
          "dataIndex": "hostName",
          "key": "hostName",
        },
        {
          "title": "人数/人数上限",
          "dataIndex": "members",
        },
        {
          "title": "游戏",
          "dataIndex": "game",
          "key": "game",
        },
        {
          "title": "活跃连接/连接总数",
          "dataIndex": "connections",
        },
        {
          "title": "运行时间",
          "dataIndex": "uptime",
        },
        {
          "title": "操作",
          "dataIndex": "operations",
        }
      ]
    }
  },
  created() {
    this.listEndpoints();
  },
  methods: {
    listEndpoints: async function() {
      try {
        const resp = await api.get("api/v1/admin/endpoints?page="+(this.page-1)+"&pageSize="+this.pageSize);
        this.endpoints = resp["endpoints"];
        this.total = resp["total"];
      }catch (_) {
        message.error("无法获取模拟器服务列表");
      }
    },

    onPageChange: function(page, pageSize) {
      this.page = page;
      this.listEndpoints();
    },

    formatMemory: function(num) {
      if (num <= 1024) return num + "Bytes";
      if (num <= 1<<20) return (num >> 10) + "KB";
      if (num <= 1<<30) return (num >> 20) + "MB";
      return (num >> 30) + "GB";
    },

    openEndpointDetail: function(endpoint) {
      this.selectedEndpoint = endpoint;
      const _this = this;
      api.get("api/v1/admin/rooms/active?id=" + endpoint["id"]).then(resp=>{
        _this.endpointDetailOpen = true;
        _this.activeRooms = resp["rooms"];
      }).catch(err=>{
        message.error("获取节点详情失败");
        _this.endpointDetailOpen = false;
      });
    },

    formatUptime: function(uptime) {
      const seconds = Math.round(parseInt(uptime)/1000);
      if (seconds <= 60) return seconds + "秒";
      if (seconds <= 3600) return Math.round(seconds/60) + "分钟";
      if (seconds <= 86400) return Math.round(seconds/3600) + "小时";
      return Math.round(seconds/86400) + "天";
    },

    closeRoom: function(id) {
      // TODO close room
    },

    deleteRoom: function(id) {
      // TODO delete room
    },
  }
}
</script>