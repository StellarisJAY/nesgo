<template>
    <a-card :bordered="false">
      <template #extra>
        <a-button v-if="joined" type="primary" @click="_=>{createRoomModalOpen = true}">新建房间</a-button>
        <a-row v-else>
          <a-col>
            <a-input v-model:value="searchInput"></a-input>
          </a-col>
          <a-col>
            <a-button type="primary" @click="searchRoom">搜索</a-button>
          </a-col>
        </a-row>
      </template>
      <a-list :grid="{ gutter: 16, xs: 1, sm: 2, md: 2, lg: 2, xl: 3, xxl: 3 }" :data-source="rooms">
        <template #renderItem="{item}">
          <a-list-item>
            <a-card :title="item.name">
              <template #extra v-if="joined">
                <a-button v-if="item.role === 0" @click="deleteRoom(item.id)" danger>删除</a-button>
                <a-button v-else @click="leaveRoom(item.id)" danger>退出</a-button>
              </template>
              <template #actions>
                <RouterLink v-if="joined" :to="'/room/' + item.id">进入</RouterLink>
                <a-button v-else type="link" @click="tryJoinRoom(item)">加入</a-button>
              </template>
              <ul style="text-align: left">
                <li>房主：{{item["hostName"]}}</li>
                <li>人数：{{item["memberCount"]}}/{{item["memberLimit"]}}</li>
              </ul>
            </a-card>
          </a-list-item>
        </template>
      </a-list>

      <a-modal v-if="joined" :open="createRoomModalOpen" title="新建房间" @cancel="_=>{createRoomModalOpen=false}">
        <template #footer>
          <a-button @click="_=>{createRoomModalOpen = false}">取消</a-button>
          <a-button type="primary" @click="createRoom()" html-type="submit">创建</a-button>
        </template>
        <a-form layout="vertical" :model="createRoomState" :label-col="{ span: 4 }">
          <a-form-item label="房间名" name="name" :rules="{required: true, message: '请输入房间名'}">
            <a-input v-model:value="createRoomState.name"></a-input>
          </a-form-item>
          <a-form-item label="私人房间" name="isPrivate">
            <a-switch v-model:checked="createRoomState.isPrivate"></a-switch>
          </a-form-item>
        </a-form>
      </a-modal>

      <a-modal v-else :open="joinRoomModalOpen" title="加入房间" @cancel="_=>{joinRoomModalOpen=false}">
        <template #footer>
          <a-button type="primary" html-type="submit" @click="enterRoom">加入</a-button>
        </template>
        <a-form layout="vertical" :modal="joinRoomFormState">
          <a-form-item label="密码" name="password">
            <a-input-password v-model:value="joinRoomFormState.password" :rules="[{required: true, message:'请输入密码'}]"></a-input-password>
          </a-form-item>
        </a-form>
      </a-modal>

    </a-card>
</template>

<script>
import { Card, Button, List, Modal, Form, Input, Switch} from 'ant-design-vue';
import { Row, Col } from "ant-design-vue";
import {message} from "ant-design-vue";
import api from "../api/request.js";
import router from "../router/index.js";
import {RouterLink} from "vue-router";

export default {
    props: {
      joined: Boolean
    },
    components: {
        ARow: Row,
        ACol: Col,
        ACard: Card,
        AButton: Button,
        AList: List,
        AListItem: List.Item,
        AModal:  Modal,
        AForm: Form,
        AFormItem: Form.Item,
        AInput: Input,
        ASwitch: Switch,
        AInputPassword: Input.Password,
    },
    data() {
        return {
           rooms: [
           ],
          createRoomModalOpen: false,
          createRoomState: {
             name: "",
              isPrivate: false,
          },
          searchInput: "",
          joinRoomModalOpen: false,
          joinRoomFormState: {
             id: 0,
             password: ""
          }
        }
    },
  created() {
    if (this.joined) {
      this.listJoinedRooms()
    }else {
      this.listAllRooms()
    }
  },
  methods: {
      listJoinedRooms() {
        api.get("api/v1/rooms/joined?page=0&pageSize=10")
            .then(resp=>{
              this.rooms = resp["rooms"];
            })
            .catch(resp=>{})
      },
      listAllRooms() {
          api.get("api/v1/rooms?page=0&pageSize=10")
              .then(resp=>{
                this.rooms = resp["rooms"];
              })
      },
      deleteRoom(id) {
        const that = this
        api.post("/room/" + id + "/delete")
            .then(resp=>{
              if (resp.status === 200) {
                message.success("删除成功")
              }
            })
            .then(_=>{
              that.listJoinedRooms(1, 10)
            })
            .catch(resp=>{})
      },
    createRoom() {
      if (this.createRoomState.name === "") {
        message.warn("请输入房间名")
        return
      }
      const that = this
      api.post("api/v1/room", {"name": this.createRoomState.name, "private": this.createRoomState.isPrivate})
          .then(resp=>{
            if (resp.status === 200) {
              message.success("新建房间成功")
              that.createRoomState = {}
              that.createRoomModalOpen = false
            }
          })
          .then(_=>{
            that.listJoinedRooms(1, 10)
          })
    },
    searchRoom() {
      api.get("/room/search?search=" + this.searchInput)
          .then(resp=>{
            this.rooms = resp.data
          })
    },
    tryJoinRoom(room) {
      api.get("/room/" + room.id + "/member")
          .then(_=>{
              return router.push("/room/" + room.id)
          })
          .catch(_=>{
              this.joinRoomFormState.id = room.id
              if (room.private) {
                this.joinRoomModalOpen = true
              }else {
                this.enterRoom()
              }
          })
    },
    enterRoom() {
      const roomId = this.joinRoomFormState.id
      const password = this.joinRoomFormState.password
      api.post("/room/" + roomId + "/join?password=" + password)
          .then(resp=>{
            this.joinRoomModalOpen = false
            return router.push("/room/" + roomId)
          })
          .catch(resp=>{
            return message.warn(resp.message)
          })
    },
    leaveRoom(id) {
      api.post("/room/" + id + "/leave")
          .then(_=>{
            message.success("退出房间成功")
            this.listJoinedRooms()
          })
    }
  }
}
</script>