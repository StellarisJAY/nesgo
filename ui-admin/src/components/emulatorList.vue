<template>
  <a-card :bordered="false">
    <a-list :grid="{ gutter: 16, xs: 1, sm: 2, md: 2, lg: 2, xl: 3, xxl: 3 }" :data-source="endpoints">
      <template #renderItem="{item}">
        <a-list-item>
          <a-card :title="item['id']">
            <template #actions>
              <a-button type="link" @click="_=>{}">详情</a-button>
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
    <a-pagination v-model:current="page" :total="total" v-model:pageSize="pageSize" @change="onPageChange" />
  </a-card>
</template>

<script>
import { Card, Button, List, Modal, Pagination} from 'ant-design-vue';
import { Row, Col } from "ant-design-vue";
import {message} from "ant-design-vue";
import api from "../api/request.js";
import router from "../router/index.js";
import {RouterLink} from "vue-router";

export default {
  props: {
  },
  components: {
    ARow: Row,
    ACol: Col,
    ACard: Card,
    AButton: Button,
    AList: List,
    AListItem: List.Item,
    AModal:  Modal,
    APagination: Pagination,
  },
  data() {
    return {
      endpoints: [
        {"id": "abc", "address": "localhost:4030"}
      ],
      page: 1,
      pageSize: 10,
      total: 0,
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
    }
  }
}
</script>