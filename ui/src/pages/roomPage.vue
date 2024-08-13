<script setup>
import {ref} from "vue";
const tourOpen = ref(true)
const refConnBtn = ref(null)
const refSelector = ref(null)
const refRestart = ref(null)
const refSaveBtn= ref(null)
const refLoadBtn = ref(null)
const refRoomBtn = ref(null)
const tourSteps = [
  {
    title: "选择游戏",
    description: "点击此处选择需要加载的游戏",
    target: ()=>refSelector.value && refSelector.value.$el,
  },
  {
    title: "连接",
    description: "点击按钮连接到模拟器",
    target: ()=>refConnBtn.value && refConnBtn.value.$el,
  },
  {
    title: "重启",
    description: "重启可以用于切换游戏，但会清除当前游戏进度，如有需要请先保存当前游戏。",
    target: ()=>refRestart.value && refRestart.value.$el,
  },
  {
    title: "保存游戏",
    description: "点击此处保存游戏进度，请注意存档数量上限。",
    target: ()=>refSaveBtn.value && refSaveBtn.value.$el,
  },
  {
    title: "读取存档",
    description: "显示存档列表，跨游戏读取存档会重启模拟器，如有需要请先保存当前游戏。",
    target: ()=>refLoadBtn.value && refLoadBtn.value.$el,
  },
  {
    title: "房间管理",
    description: "点击此处弹出房间面板，房主可通过此面板修改房间信息以及玩家权限。",
    target: ()=>refRoomBtn.value && refRoomBtn.value.$el,
  }
]
</script>

<template>
  <a-row style="height: 100vh; background-color: grey">
    <!--left side buttons-->
    <a-col :span="7">
      <a-row style="height: 30%; margin-top: 10%">
        <a-col :span="8" :offset="8">
          <a-button class="control-btn" id="button-up" disabled><ArrowUpOutlined /></a-button>
        </a-col>
      </a-row>
      <a-row style="height: 30%">
        <a-col :span="8">
          <a-button class="control-btn" id="button-left" disabled><ArrowLeftOutlined /></a-button>
        </a-col>
        <a-col :span="8" :offset="8">
          <a-button class="control-btn" id="button-right" disabled><ArrowRightOutlined /></a-button>
        </a-col>
      </a-row>
      <a-row style="height: 30%">
        <a-col :span="8" :offset="8">
          <a-button class="control-btn" id="button-down" disabled><ArrowDownOutlined/></a-button>
        </a-col>
      </a-row>
    </a-col>
    <!--video screen and toolbar-->
    <a-col :span="10" style="height: 100vh">
      <a-card style="height: 100%">
        <a-row>
          <a-col :span="6">
            <a-button ref="refSaveBtn" type="primary" :disabled="saveBtnDisabled" @click="saveGame" style="width: 90%">保存</a-button>
          </a-col>
          <a-col :span="6">
            <a-button ref="refLoadBtn" type="primary" :disabled="loadBtnDisabled" @click="openSavedGamesDrawer" style="width: 90%">读档</a-button>
          </a-col>
          <a-col :span="6">
            <a-button type="primary" :disabled="chatBtnDisabled" @click="_=>{setChatModal(true)}" style="width: 90%">聊天</a-button>
          </a-col>
          <a-col :span="6">
            <a-button ref="refRoomBtn" type="primary" @click="openRoomMemberDrawer" style="width: 90%">房间</a-button>
          </a-col>
        </a-row>
        <a-row style="height: 80%">
          <video id="video" playsinline webkit-playsinline="true"></video>
        </a-row>
        <a-row>
          <a-col :span="6">
            <a-button ref="refConnBtn" style="width: 90%;" type="primary" @click="connect" :disabled="connectBtnDisabled">连接</a-button>
          </a-col>
          <a-col :span="6">
            <a-button ref="refRestart" style="width: 90%;" type="primary" :disabled="restartBtnDisabled" @click="restart">重启</a-button>
          </a-col>
          <a-col>
            <a-select
                ref="refSelector"
                v-model:value="selectedGame"
                :options="configs.existingGames"
            ></a-select>
          </a-col>
        </a-row>
      </a-card>
    </a-col>
    <!--right side buttons-->
    <a-col :span="7">
      <a-row style="height: 30%; margin-top: 10%">
        <a-col :span="8">
          <a-button class="control-btn" id="button-start" disabled>START</a-button>
        </a-col>
        <a-col :span="8" :offset="8">
          <a-button class="control-btn" id="button-a" disabled>A</a-button>
        </a-col>
      </a-row>
      <a-row style="height: 30%; margin-top: 60%">
        <a-col :span="8">
          <a-button class="control-btn" id="button-select" disabled>SELECT</a-button>
        </a-col>
        <a-col :span="8" :offset="8">
          <a-button class="control-btn" id="button-b" disabled>B</a-button>
        </a-col>
      </a-row>
    </a-col>
    <!--room member list-->
    <a-drawer
        v-model:open="membersDrawerOpen"
        :root-style="{ color: 'blue' }"
        title="房间信息"
        placement="right"
    >
      <a-form v-if="memberSelf.role === RoleNameHost">
        <a-form-item label="私人">
          <a-switch v-model:checked="fullRoomInfo.private" @change="alterRoomPrivacy"></a-switch>
        </a-form-item>
        <a-form-item label="密码" v-if="fullRoomInfo.private">
          <a-input-password readonly :value="fullRoomInfo.password"></a-input-password>
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
    </a-drawer>
    <!--saved games-->
    <a-drawer size="default" title="保存游戏" placement="right" v-model:open="savedGameOpen">
      <a-list item-layout="vertical" :data-source="savedGames">
        <template #renderItem="{item}">
          <a-list-item>
            <a-descriptions :column="1">
              <a-descriptions-item label="游戏">{{item["game"]}}</a-descriptions-item>
              <a-descriptions-item label="时间">{{item["createdAt"]}}</a-descriptions-item>
            </a-descriptions>
            <a-button type="primary" @click="loadSavedGame(item.id)">加载</a-button>
            <a-button danger @click="deleteSavedGame(item.id)">删除</a-button>
          </a-list-item>
        </template>
      </a-list>
    </a-drawer>

    <a-modal title="聊天" v-model:open="chatModalOpen" @cancel="_=>{setChatModal(false)}">
      <template #footer>
        <a-button @click="_=>{setChatModal(false)}">取消</a-button>
        <a-button type="primary" @click="sendChatMessage">发送</a-button>
      </template>
      <a-input placeholder="请输入消息..." v-model:value="chatMessage"></a-input>
    </a-modal>
    <a-tour :steps="tourSteps" :open="tourOpen" @close="_=>{tourOpen=false}"></a-tour>
  </a-row>
  <p id="rtt">RTT: {{rtt}}ms</p>
</template>

<script>
import api from "../api/request.js";
import { Row, Col } from "ant-design-vue";
import {CrownTwoTone} from '@ant-design/icons-vue';
import {Card, Button, Drawer, List, Descriptions, RadioGroup, Radio, Select, Checkbox, InputPassword, Switch} from "ant-design-vue";
import {message} from "ant-design-vue";
import {Form, FormItem, Modal, Input} from "ant-design-vue";
import tokenStorage from "../api/token.js";
import router from "../router/index.js";
import {ArrowUpOutlined, ArrowDownOutlined, ArrowLeftOutlined, ArrowRightOutlined, SaveOutlined} from "@ant-design/icons-vue"
import {notification, Tour} from "ant-design-vue";

// const MessageSDPOffer = 0
// const MessageSDPAnswer = 1
// const MessageICECandidate = 2
const MessageGameButtonPressed = 0
const MessageGameButtonReleased = 1
// const MessageTurnServerInfo = 5
// const MessageChat = 6
// const MessagePing = 7
// const MessagePong = 8

// const RoleHost = 0
// const RoleGamer = 1
// const RoleObserver = 2

const RoleNameHost = "Host";
const RoleNamePlayer = "Player";
const RoleNameObserver = "Observer";

export default {
  components: {
    ARow: Row,
    ACol: Col,
    ACard: Card,
    AButton: Button,
    ADrawer: Drawer,
    AList: List,
    AListItem: List.Item,
    ARadio: Radio,
    ARadioGroup: RadioGroup,
    CrownTwoTone: CrownTwoTone,
    ArrowUpOutlined,
    ArrowDownOutlined,
    ArrowLeftOutlined,
    ArrowRightOutlined,
    ASelect: Select,
    ADescriptions: Descriptions,
    ADescriptionsItem: Descriptions.Item,
    SaveOutlined,
    ACheckbox: Checkbox,
    AInputPassword: InputPassword,
    ASwitch: Switch,
    AForm: Form,
    AFormItem: FormItem,
    AModal: Modal,
    AInput: Input,
    ATour: Tour,
  },
    data() {
        return {
            members: [],
            membersDrawerOpen: false,
            memberSelf: {
              role: 3,
            },
            rtcSession: {},
            connectBtnDisabled: false,
            saveBtnDisabled: true,
            loadBtnDisabled: true,
            restartBtnDisabled: true,
            chatBtnDisabled: true,
            selectedGame: "",
            configs: {
              controlButtonMapping: {
                "button-up": "Up",
                "button-down": "Down",
                "button-left": "Left",
                "button-right": "Right",
                "button-a": "A",
                "button-b": "B",
                "button-select": "Select",
                "button-start": "Start",
              },
              keyboardMapping: {
                "KeyA": "Left",
                "KeyD": "Right",
                "KeyW": "Up",
                "KeyS": "Down",
                "KeyJ": "B",
                "Space": "A",
                "Enter": "Start",
                "Tab": "Select",
              },
              existingGames: [],
            },
            savedGameOpen: false,
            savedGames: [],
            p1p2Options: [
              {value: "1", label: "P1"},
              {value: "2", label: "P2"},
            ],
            fullRoomInfo: {},
            chatModalOpen: false,
            chatMessage: "",
            pingInterval: 0,
            rtt: 0,
            iceCandidates: [],
        }
    },
  created() {
      this.roomId = this.$route["params"]["roomId"];
      this.getMemberSelf();
      this.listGames();
  },
  unmounted() {
    this.rtcSession.pc.close()
  },
  methods: {
      openRoomMemberDrawer() {
          this.listRoomMembers().then(_=>this.membersDrawerOpen=true);
      },

      getMemberSelf: async function() {
        const resp = await api.get("api/v1/member/"+this.roomId);
        const member = resp["member"];
        this.memberSelf = member;
        if (member.role === RoleNameHost) {
          this.fullRoomInfo = await api.get("api/v1/room/"+this.roomId);
        }
        await this.listRoomMembers();
      },

      listRoomMembers: async function() {
        const resp = await api.get("api/v1/members?roomId=" + this.roomId);
        this.members = resp["members"];
      },

      listGames: async function() {
        const resp = await api.get("api/v1/games");
        this.games = resp["games"];
        const existingGames = [];
        this.games.forEach(game=>{
          existingGames.push({"label": game["name"], "value": game["name"]});
        });
        this.configs.existingGames = existingGames;
        this.selectedGame = existingGames[0].value;
      },

      connect() {
        this.connectBtnDisabled = true
        this.openConnection();
      },

      openConnection: async function() {
        const _this = this;
        const roomId = this.roomId;
        const data = await api.post("api/v1/game/connection", {
          "roomId": roomId,
          "game": this.selectedGame,
        });
        // TODO TURN relay server
        const pc = new RTCPeerConnection({
          iceServers: [
            {
              urls: "stun:stun.l.google.com"
            },
          ],
          iceTransportPolicy: "all",
        });
        // on remote track
        pc.ontrack = ev=>{
          console.log("on track: ", ev);
          if (ev.track.kind === "video") {
            document.getElementById("video").srcObject = ev.streams[0]
            document.getElementById("video").autoplay = true
            document.getElementById("video").controls = true
          }
        };

        // 发送answer之前的candidate，避免远端没有收到answer导致无法这是candidate
        pc.onicecandidate = ev=>{
          if (ev.candidate) {
            _this.iceCandidates.push(ev.candidate);
          }
        };

        pc.onconnectionstatechange = ev=>this.onPeerConnStateChange(ev);

        const rtcSession = {
          roomId: roomId,
          pc: pc,
          dataChannel: null,
        }
        await pc.setRemoteDescription({
          type: "offer",
          sdp: data["sdpOffer"],
        });
        const answer = await pc.createAnswer();
        await pc.setLocalDescription(answer);
        await api.post("api/v1/game/sdp", {
          "roomId": roomId,
          "sdpAnswer": answer.sdp,
        });
        // 发送answer之前的candidate，避免远端没有收到answer导致无法这是candidate
        this.iceCandidates.forEach(candidate=>{
          const s = JSON.stringify(candidate);
          api.post("api/v1/game/ice", {
            "roomId": roomId,
            "candidate": s,
          });
        });
        // 发送answer之后的candidate直接发送给远端
        pc.onicecandidate = ev=>{
          if (ev.candidate) {
            const s = JSON.stringify(ev.candidate);
            api.post("api/v1/game/ice", {
              "roomId": roomId,
              "candidate": s,
            });
          }
        }
        // data channel
        pc.ondatachannel = ev=>{
          console.log("on datachannel: ", ev);
          rtcSession.dataChannel = ev.channel;
          ev.channel.onopen = _=>console.log("datachannel open");
          // TODO handle data channel message
          ev.channel.onmessage = msg=>{
            console.log(msg.data);
          };
        };

        this.rtcSession = rtcSession;
      },

      onPeerConnStateChange(_) {
        const pc = this.rtcSession.pc
        console.log("peer conn state: " + pc.connectionState)
        switch (pc.connectionState) {
          case "connected":
            this.onConnected()
            break
          case "disconnected":
            this.onDisconnected()
            break
          default:
            break
        }
      },
      onConnected() {
        message.success("连接成功");
        if (this.memberSelf["role"] !== RoleNameObserver) {
          this.setKeyboardControl(true);
          this.initControlButtons();
        }
        this.saveBtnDisabled = this.memberSelf["role"] !== RoleNameHost;
        this.loadBtnDisabled = false;
        this.restartBtnDisabled = this.memberSelf["role"] !== RoleNameHost;
      },
    onDisconnected() {
      message.warn("连接断开");
      this.disableControlButtons()
      this.rtcSession.pc.close()
      this.connectBtnDisabled = false
    },
    sendAction(code, pressed) {
        // TODO send control message
      const msg = JSON.stringify({
        "from": 0,
        "to": 0,
        "type": pressed,
        "data": code,
      });
      this.rtcSession.dataChannel.send(msg);
    },

    onRoleRatioChange(ev, member) {
        member.role = ev.target.value;
        console.log("change member role: " + member.role);
    },
    kickMember(member) {
        console.log("kick member: " + member);
    },
    initControlButtons() {
        const mapping = this.configs.controlButtonMapping
        for (const k in mapping) {
          const button = document.getElementById(k);
          const keyCode = mapping[k]
          button.disabled = false
          button.addEventListener("touchstart", _=>this.sendAction(keyCode, MessageGameButtonPressed))
          button.addEventListener("touchend", _=>this.sendAction(keyCode, MessageGameButtonReleased))
          button.addEventListener("mousedown", _=>this.sendAction(keyCode, MessageGameButtonPressed))
          button.addEventListener("mouseup", _=>this.sendAction(keyCode, MessageGameButtonReleased))
        }
    },
    disableControlButtons() {
      const mapping = this.configs.controlButtonMapping
      for (const k in mapping) {
        const button = document.getElementById(k)
        button.disabled = true
      }
    },
    getSavedGames() {
      // TODO list saved games
    },
    openSavedGamesDrawer() {
        this.getSavedGames().then(_=>this.savedGameOpen=true)
    },
    saveGame() {
        // TODO save game
    },
    loadSavedGame(id) {
        // TODO load save
    },
    deleteSavedGame(id) {
        // TODO delete save
    },
    restart() {
        // TODO restart emulator
    },

    onP1P2Change(ev, m, which) {
      if (m["role"] === RoleNameObserver) {
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
    setChatModal(open) {
      this.setKeyboardControl(!open)
      this.chatModalOpen = open
      if (!open) {
        this.chatMessage = ""
      }
    },
    sendChatMessage() {
      // TODO send chat message
    },

    setKeyboardControl(enabled) {
        if (enabled) {
          window.onkeydown = ev=> {
            const button = this.configs.keyboardMapping[ev.code];
            if (button) {
              this.sendAction(button, MessageGameButtonPressed)
            }
          }

          window.onkeyup = ev=> {
            const button = this.configs.keyboardMapping[ev.code];
            if (button) {
              this.sendAction(button, MessageGameButtonReleased)
            }
          }
        }else {
          window.onkeydown = _=>{}
          window.onkeyup = _=>{}
        }
    },

    ping() {
        const timestamp = new Date().getTime()
        // TODO send ping message
    }
  }
}
</script>

<style scoped>
#video {
  width: 100%;
  background-color: black;
}
.control-btn {
  width: 100%;
  height: 100%;
}
#rtt {
  position: absolute;
  right: 0;
  top: 0;
}
</style>