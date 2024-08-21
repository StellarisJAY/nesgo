<template>
  <a-card :bordered="false">
    <!--游戏列表-->
    <a-table :data-source="games" :columns="columns" :pagination="false">
      <template #bodyCell="{column, text ,record}">
        <template v-if="column.dataIndex === 'operation'">
          <a-button type="primary" danger @click="deleteGame(record['name'])">删除</a-button>
        </template>
        <template v-else-if="column.dataIndex === 'size'">
          {{formatMemory(text)}}
        </template>
        <template v-else-if="column.dataIndex === 'uploadTime'">
          {{new Date(parseInt(record["uploadTime"])).toLocaleString()}}
        </template>
      </template>
    </a-table>

    <!--分页-->
    <a-pagination v-model:current="page" :total="total" v-model:pageSize="pageSize" @change="onPageChange" />

    <!--上传游戏-->
    <a-modal :open="uploadGameModalOpen" title="上传游戏" @cancel="cancelUpload">
      <template #footer>
        <a-button type="primary" @click="cancelUpload">取消</a-button>
        <a-button type="primary" @click="uploadGame">上传</a-button>
      </template>
      <a-upload-dragger v-model:file-list="fileList" :multiple="false" :before-upload="_ => false">
        <p class="ant-upload-drag-icon">
          <inbox-outlined></inbox-outlined>
        </p>
        <p class="ant-upload-text">点击或拖拽文件到该区域</p>
      </a-upload-dragger>
    </a-modal>

    <template #extra>
      <a-button type="primary" @click="_=>{uploadGameModalOpen = true;}">上传游戏</a-button>
    </template>
  </a-card>
</template>

<script>
import { Card, Button, List, Modal, Pagination, Table, UploadDragger} from 'ant-design-vue';
import {InboxOutlined} from '@ant-design/icons-vue';
import { Row, Col } from "ant-design-vue";
import {message} from "ant-design-vue";
import api from "../api/request.js";

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
    AUploadDragger: UploadDragger,
    InboxOutlined,
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
        {
          "title": "操作",
          "dataIndex": "operation",
        }
      ],
      page: 1,
      pageSize: 10,
      total: 0,

      uploadGameModalOpen: false,
      fileList: [],
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
        this.total = resp["total"];
      }catch (_) {
        message.warn("获取失败");
      }
    },
    onPageChange: function(page, pageSize) {
      this.page = page;
      this.listGames();
    },
    formatMemory: function(num) {
      if (num <= 1024) return num + "Bytes";
      if (num <= 1<<20) return (num >> 10) + "KB";
      if (num <= 1<<30) return (num >> 20) + "MB";
      return (num >> 30) + "GB";
    },

    uploadGame: function() {
      const _this = this;
      let formData = new FormData();
      formData.append("data", this.fileList[0].originFileObj);
      api.upload("api/v1/admin/game/upload", formData).then(_=>{
        _this.fileList = [];
        message.success("上传成功");
        _this.listGames();
        _this.uploadGameModalOpen = false;
      }).catch(err=>{
        message.error("上传失败，请检查文件格式");
      });
    },

    cancelUpload: function() {
      this.fileList = [];
      this.uploadGameModalOpen = false;
    },

    deleteGame: function(name) {
      const _this = this;
      api.delete("api/v1/admin/games?games[]=" + name).then(_=>{
        message.success("删除成功");
        _this.listGames();
      }).catch(_=>{
        message.error("删除失败");
      });
    },
  }
}
</script>