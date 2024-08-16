<template>
  <div>
    <a-form v-if="memberSelf['role'] === RoleNameHost">
      <a-form-item label="私人">
        <a-switch v-model:checked="fullRoomInfo['private']" @change="alterRoomPrivacy"></a-switch>
      </a-form-item>
      <a-form-item label="密码" v-if="fullRoomInfo['private']">
        <a-input-password readonly :value="fullRoomInfo['password']"></a-input-password>
      </a-form-item>
    </a-form>
    <a-list item-layout="vertical" :data-source="members">
      <template #renderItem="{item}">
        <a-list-item>
          <a-row>
            <a-col :span="8"><CrownTwoTone v-if="item.role===RoleNameHost" />{{item.name}}</a-col>
            <a-col :span="10">
              <a-checkbox :disabled="memberSelf.role!==RoleNameHost||rtcSession.pc === undefined||rtcSession.pc.connectionState !== 'connected' "
                          v-model:checked="item['player1']"
                          @change="ev=>{onP1P2Change(ev, item, 1)}">P1</a-checkbox>
              <a-checkbox :disabled="memberSelf.role!==RoleNameHost||rtcSession.pc === undefined||rtcSession.pc.connectionState !== 'connected' "
                          v-model:checked="item['player2']"
                          @change="ev=>{onP1P2Change(ev, item, 2)}">P2</a-checkbox>
              <a-radio-group v-model:value="item.role" :disabled="memberSelf.role!==RoleNameHost || item.role===RoleNameHost"
                             @change="ev=>{onRoleRatioChange(ev, item)}">
                <a-radio :value="RoleNamePlayer">玩家</a-radio>
                <a-radio :value="RoleNameObserver">观战</a-radio>
              </a-radio-group>
            </a-col>
            <a-col :span="6">
              <a-button type="primary" :hidden="memberSelf.role!==RoleNameHost"
                        :disabled="item.role===RoleNameHost || memberSelf.role!==RoleNameHost"
                        @click="kickMember(item)">踢出</a-button>
            </a-col>
          </a-row>
        </a-list-item>
      </template>
    </a-list>
  </div>
</template>

<script>
import { List, Button, Checkbox, Drawer, Form, Input, Switch, Radio, RadioGroup} from 'ant-design-vue';
import { Row, Col } from "ant-design-vue";
import {message} from "ant-design-vue";
import api from "../api/request.js";
import {CrownTwoTone} from "@ant-design/icons-vue"

export default {
  props: {
    memberSelf: Object,
    fullRoomInfo: Object,
    rtcSession: Object,
    roomId: String,
  },
  components: {
    AButton: Button,
    ADrawer: Drawer,
    AForm: Form,
    AInput: Input,
    ARow: Row,
    ACol: Col,
    ASwitch: Switch,
    ACheckbox: Checkbox,
    AList: List,
    AListItem: List.Item,
    ARadio: Radio,
    ARadioGroup: RadioGroup,
    CrownTwoTone,
    AInputPassword: Input.Password,
    AFormItem: Form.Item,
  },
  data() {
    return {
      members: [],
      RoleNameHost: "Host",
      RoleNamePlayer: "Player",
      RoleNameObserver: "Observer",
    }
  },
  created() {
    this.listRoomMembers();
    addEventListener("memberDrawerOpen", _=>this.listRoomMembers());
  },
  methods: {
    listRoomMembers: async function() {
      const resp = await api.get("api/v1/members?roomId=" + this.roomId);
      this.members = resp["members"];
    },
    onRoleRatioChange(ev, member) {
      member.role = ev.target.value;
      const _this = this;
      api.put("api/v1/member/role", {
        "roomId": this.roomId,
        "userId": member["userId"],
        "role": member["role"],
      }).then(_=>{
        message.success("操作成功");
        _this.listRoomMembers();
      }).catch(_=>{
        message.warn("操作失败");
      });
    },
    kickMember(member) {
      const _this = this;
      api.delete("api/v1/member?roomId="+this.roomId+"&userId="+member["userId"]).then(_=>{
        message.success("操作成功");
        _this.listRoomMembers();
      }).catch(err=>{
        message.error("操作失败");
      });
    },

    onP1P2Change(ev, m, which) {
      if (m["role"] === this.RoleNameObserver) {
        message.error("无法修改观战玩家的控制");
        return;
      }
      const _this = this;
      api.post("api/v1/game/controller", {
        "roomId": this.roomId,
        "playerId": m["userId"],
        "controllerId": ev.target.checked ? which : -1,
      }).then(_=>{
        _this.listRoomMembers();
        message.success("修改成功");
      });
    },
    alterRoomPrivacy() {
      // TODO update room
    },
  }
}
</script>