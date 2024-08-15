<template>
  <a-list item-layout="vertical" :data-source="saves">
    <template #renderItem="{item}">
      <a-list-item>
        <a-descriptions :column="1">
          <a-descriptions-item label="游戏">{{item["game"]}}</a-descriptions-item>
          <a-descriptions-item label="时间">{{new Date(parseInt(item["createTime"])).toLocaleString()}}</a-descriptions-item>
        </a-descriptions>
        <a-button type="primary" @click="loadSavedGame(item.id)">加载</a-button>
        <a-button danger @click="deleteSavedGame(item.id)">删除</a-button>
      </a-list-item>
    </template>
    <a-pagination :page="page" :pageSize="pageSize" :total="total" @change="onPaginationChange"></a-pagination>
  </a-list>
</template>

<script>
import { List, Button, Descriptions, Pagination} from 'ant-design-vue';
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
  },
  data() {
    return {
      saves: [],
      page: 1,
      pageSize: 10,
      total: 0,
    }
  },
  created() {
    this.listSaves();
  },
  methods: {
    listSaves: function() {
      const _this = this;
      api.get("api/v1/game/saves?roomId="+this.roomId+"&page="+(this.page-1)+"&pageSize="+this.pageSize).then(resp=>{
        _this.saves = resp["saves"];
        _this.total = resp["total"];
      }).catch(_=>{
        message.error("获取存档列表失败");
      })
    },
    onPaginationChange: function(page, pageSize) {
      this.page = page;
      this.listSaves();
    },
    loadSavedGame(id) {
      // TODO load save
    },
    deleteSavedGame(id) {
      // TODO delete save
    },
  }
}
</script>