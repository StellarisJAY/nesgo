<template>
  <a-card :bordered="false">
    <a-table :data-source="games" :columns="columns"></a-table>
    <a-pagination v-model:current="page" :total="total" v-model:pageSize="pageSize" @change="onPageChange" />
  </a-card>
</template>

<script>
import { Card, Button, List, Modal, Pagination, Table} from 'ant-design-vue';
import { Row, Col } from "ant-design-vue";
import {message} from "ant-design-vue";
import api from "../api/request.js";
// import router from "../router/index.js";
// import {RouterLink} from "vue-router";

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
    ATable: Table,
  },
  data() {
    return {
      games: [
        // {"name": "SuperMario.nes", "mapper": "mapper0", "mirroring": "vertical", "size": 11231, "uploadTime": "2024-08-10"}
      ],
      columns: [
        {
          "title": "名称",
          "dataIndex": "name",
          "key": "name",
        },
        {
          "title": "卡带格式",
          "dataIndex": "mapper",
          "key": "mapper",
        },
        {
          "title": "画面映射",
          "dataIndex": "mirroring",
          "key": "mirroring",
        },
        {
          "title": "文件大小",
          "dataIndex": "size",
          "key": "size",
        },
        {
          "title": "上传时间",
          "dataIndex": "uploadTime",
          "key": "uploadTime",
        },
      ],
      page: 1,
      pageSize: 10,
      total: 0,
    }
  },
  created() {
    this.listGames();
  },
  methods: {
    listGames: async function() {
      try {
        const resp = await api.get("api/v1/admin/games?page=" + (this.page - 1) + "&pageSize=" + this.pageSize);
        this.games = resp["games"];
      }catch (_) {
        message.warn("获取失败");
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