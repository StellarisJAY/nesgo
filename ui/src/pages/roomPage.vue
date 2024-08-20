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
          <a-col :span="6"></a-col>
          <a-col :span="6">
            <a-button style="width: 90%" type="primary" @click="openSettingDrawer">设置</a-button>
          </a-col>
        </a-row>
        <a-row>
          <a-col :span="24">
            <a-select
                ref="refSelector"
                v-model:value="selectedGame"
                :options="configs.existingGames"
                style="width: 100%"
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
    <a-drawer v-model:open="membersDrawerOpen" placement="right" title="房间信息" size="default">
      <RoomInfoDrawer :member-self="memberSelf" :rtc-session="rtcSession" :full-room-info="fullRoomInfo" :room-id="roomId"></RoomInfoDrawer>
    </a-drawer>
    <!--saved games-->
    <a-drawer size="default" title="保存游戏" placement="right" v-model:open="savedGameOpen">
     <SaveList :room-id="roomId"></SaveList>
    </a-drawer>
    <!--chat modal-->
    <a-modal title="聊天" v-model:open="chatModalOpen" @cancel="_=>{setChatModal(false)}">
      <template #footer>
        <a-button @click="_=>{setChatModal(false)}">取消</a-button>
        <a-button type="primary" @click="sendChatMessage">发送</a-button>
      </template>
      <a-input placeholder="请输入消息..." v-model:value="chatMessage"></a-input>
    </a-modal>
    <!--settings-->
    <a-drawer v-model:open="settingDrawerOpen" placement="right" title="设置" size="default">
      <p>提示：点击按钮取消绑定，点击‘+’后按下键盘按键添加绑定</p>
      <a-row>
        <a-col :span="12">
          <a-select :options="configs.bindingOptions" v-model:value="configs.selectedBindingKey" @change="onKeyboardBindingSelectChange"></a-select>
        </a-col>
        <a-col :span="6"></a-col>
        <a-col :span="6">
          <a-button type="primary" @click="setKeyboardBindingEnabled">启用</a-button>
        </a-col>
      </a-row>
      <!--按键设置列表-->
      <a-table :data-source="configs.selectedBinding['bindings']"
               :columns="configs.bindingTableColumns"
               :pagination="false">
        <template #bodyCell="{column, record}">
          <template v-if="column.dataIndex === 'emulatorKey'">
            {{record.emulatorKey}}
          </template>
          <template v-else-if="column.dataIndex==='keyboardKey'">
            <KeyboardKeyPicker :limit="1" :buttons="record.buttons"></KeyboardKeyPicker>
          </template>
        </template>
      </a-table>
      <!--保存按键设置按钮-->
      <a-row>
        <a-col :span="18"></a-col>
        <a-col :span="6">
          <a-button type="primary" @click="updateBinding" :disabled="configs.selectedBinding.id==='0'">保存修改</a-button>
        </a-col>
      </a-row>
      <a-form>
        <a-form-item label="显示状态数据">
          <a-switch v-model:checked="configs.showStats"></a-switch>
        </a-form-item>
      </a-form>
    </a-drawer>
    <a-tour :steps="tourSteps" :open="tourOpen" @close="_=>{tourOpen=false}"></a-tour>
  </a-row>
  <p id="stats" v-if="configs.showStats">RTT:{{stats.rtt }}ms FPS:{{stats.fps}}</p>
</template>

<script>
import api from "../api/request.js";
import globalConfigs from "../api/const.js";
import { Row, Col } from "ant-design-vue";
import {CrownTwoTone} from '@ant-design/icons-vue';
import {Card, Button, Drawer, List, Descriptions, RadioGroup, Radio, Select, Checkbox, InputPassword, Switch, Table, notification} from "ant-design-vue";
import {message} from "ant-design-vue";
import {Form, FormItem, Modal, Input} from "ant-design-vue";
import router from "../router/index.js";
import {ArrowUpOutlined, ArrowDownOutlined, ArrowLeftOutlined, ArrowRightOutlined, SaveOutlined} from "@ant-design/icons-vue"
import {Tour} from "ant-design-vue";
import RoomInfoDrawer from "../components/roomInfoDrawer.vue";
import SaveList from "../components/saveList.vue";
import KeyboardKeyPicker from "../components/keyboardKeyPicker.vue";

const MessageGameButtonPressed = 0
const MessageGameButtonReleased = 1
const MessageChat = 2;
const MessagePing = 14;

const RoleNameHost = "Host";
const RoleNamePlayer = "Player";
const RoleNameObserver = "Observer";

const defaultBinding = {
      "id": "0",
      "name": "默认绑定",
      "bindings": [
        {
          "emulatorKey": "Left",
          "emulatorKeyTranslated": "Left",
          "buttons": ["KeyA"],
          "keyboardKey": "KeyA",
        },
        {
          "emulatorKey": "Right",
          "emulatorKeyTranslated": "Right",
          "buttons": ["KeyD"],
          "keyboardKey": "KeyD",
        },
        {
          "emulatorKey": "Up",
          "emulatorKeyTranslated": "Up",
          "buttons": ["KeyW"],
          "keyboardKey": "KeyW",
        },
        {
          "emulatorKey": "Down",
          "emulatorKeyTranslated": "Down",
          "buttons": ["KeyS"],
          "keyboardKey": "KeyS",
        },
        {
          "emulatorKey": "A",
          "emulatorKeyTranslated": "A",
          "buttons": ["Space"],
          "keyboardKey": "Space",
        },
        {
          "emulatorKey": "B",
          "emulatorKeyTranslated": "B",
          "buttons": ["KeyJ"],
          "keyboardKey": "KeyJ",
        },
        {
          "emulatorKey": "Start",
          "emulatorKeyTranslated": "Start",
          "buttons": ["Enter"],
          "keyboardKey": "Enter",
        },
        {
          "emulatorKey": "Select",
          "emulatorKeyTranslated": "Select",
          "buttons": ["Tab"],
          "keyboardKey": "Tab",
        },
      ]
    };

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
    RoomInfoDrawer,
    SaveList,
    ATable: Table,
    KeyboardKeyPicker,
  },
    data() {
        return {
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
              bindingTableColumns: [
                {
                  "title": "键盘按键",
                  "dataIndex": "keyboardKey",
                  "key": "keyboardKey",
                },
                {
                  "title": "模拟器按键",
                  "dataIndex": "emulatorKey",
                  "key": "emulatorKey",
                }
              ],
              userBindings: [],
              selectedBinding: defaultBinding,
              selectedBindingKey: "",
              bindingOptions: [],
              existingGames: [],
              showStats: false,
            },
            savedGameOpen: false,
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

            settingDrawerOpen: false,

            stats: {
              rtt: 0,
              fps: 0,
            }
        }
    },
  created() {
      this.roomId = this.$route["params"]["roomId"];
      this.getMemberSelf();
      this.listGames();
  },
  unmounted() {
    if (this.rtcSession && this.rtcSession.pc) this.rtcSession.pc.close();
  },
  methods: {
      openRoomMemberDrawer() {
          this.membersDrawerOpen = true;
          dispatchEvent(new Event("memberDrawerOpen"));
      },

      getMemberSelf: async function() {
        try {
          const resp = await api.get("api/v1/member/"+this.roomId);
          const member = resp["member"];
          this.memberSelf = member;
          if (member.role === RoleNameHost) {
            this.fullRoomInfo = await api.get("api/v1/room/"+this.roomId);
          }
        }catch (_) {
          message.warn("无法访问房间");
          router.push("/home");
        }
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
        let data;
        try {
          data = await api.post("api/v1/game/connection", {
            "roomId": roomId,
            "game": this.selectedGame,
          });
        }catch (errResp) {
          message.warn("连接失败，请重试");
          this.connectBtnDisabled = false;
          return;
        }

        const pc = new RTCPeerConnection({
          iceServers: [
            {
              urls: globalConfigs.StunServer,
            },
            {
              urls: globalConfigs.TurnServer.Host,
              username: globalConfigs.TurnServer.Username,
              credential: globalConfigs.TurnServer.Password,
            }
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
        try {
          await api.post("api/v1/game/sdp", {
            "roomId": roomId,
            "sdpAnswer": answer.sdp,
          });
        }catch (errResp) {
          message.warn("连接失败，请重试");
          this.connectBtnDisabled = false;
          return;
        }

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
            }).then(_=>{
              return api.get("api/v1/ice/candidates?roomId=" + this.roomId);
            }).then(resp=>{
              resp["candidates"].forEach(candidate=>{
                const c = JSON.parse(candidate);
                console.log("remote candidate: ", c);
                pc.addIceCandidate(c);
              });
            })
          }
        }
        // data channel
        pc.ondatachannel = ev=>{
          rtcSession.dataChannel = ev.channel;
          ev.channel.onopen = _=>_this.onDataChannelOpen();
          ev.channel.onmessage = msg=>_this.onDataChannelMsg(msg);
          ev.channel.onclose = _=>_this.onDataChannelClose();
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
      this.disableControlButtons();
      this.rtcSession.pc.close();
      this.saveBtnDisabled = true;
      this.restartBtnDisabled = true;
      this.loadBtnDisabled = true;
      this.chatBtnDisabled = true;
      this.getMemberSelf().then(_=>{
        this.connectBtnDisabled = false;
      }).catch(_=>{
        message.warn("无法访问该房间");
        router.back();
      })
    },
    sendAction(code, pressed) {
      const msg = JSON.stringify({
        "type": pressed,
        "data": code,
      });
      this.rtcSession.dataChannel.send(msg);
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

    openSavedGamesDrawer() {
        this.savedGameOpen = true;
        dispatchEvent(new Event("saveListOpen"));
    },
    saveGame() {
        const _this = this;
        this.saveBtnDisabled = true;
        api.post("api/v1/game/save", {
          "roomId": this.roomId,
        }).then(_=>{
          return message.success("保存成功");
        }).then(_=>{
          _this.saveBtnDisabled = false;
        }).catch(_=>{
          message.error("保存失败");
          _this.saveBtnDisabled = false;
        })
    },

    restart() {
      this.restartBtnDisabled = true;
      const _this = this;
      api.post("api/v1/game/restart", {
        "roomId": this.roomId,
        "game":   this.selectedGame,
      }).then(_=>{
        return message.success("重启成功")
      }).then(_=>{
        _this.restartBtnDisabled = false;
      }).catch(_=>{
        message.error("重启失败");
        _this.restartBtnDisabled = false;
      })
    },

    setChatModal(open) {
      this.setKeyboardControl(!open)
      this.chatModalOpen = open
      if (!open) {
        this.chatMessage = ""
      }
    },
    sendChatMessage() {
      const timestamp = new Date().getTime();
      if (this.rtcSession && this.rtcSession.pc) {
        const pingMsg = {
          "type": MessageChat,
          "timestamp": timestamp,
          "data": this.chatMessage,
        };
        this.rtcSession.dataChannel.send(JSON.stringify(pingMsg));
      }
      this.setChatModal(false);
    },

    setKeyboardControl(enabled) {
        if (enabled) {
          const _this = this;
          window.onkeydown = ev=>{
            const button = _this.configs.selectedBinding.bindings.find(item=>item.buttons[0] === ev.code);
            if (button) {
              _this.sendAction(button.emulatorKey, MessageGameButtonPressed);
            }
          };

          window.onkeyup = ev=>{
            const button = _this.configs.selectedBinding.bindings.find(item=>item.buttons[0] === ev.code);
            if (button) {
              _this.sendAction(button.emulatorKey, MessageGameButtonReleased);
            }
          };
        }else {
          window.onkeydown = _=>{}
          window.onkeyup = _=>{}
        }
    },

    ping() {
        const timestamp = new Date().getTime();
        if (this.rtcSession && this.rtcSession.pc) {
          const pingMsg = {
            "type": MessagePing,
            "timestamp": timestamp,
          };
          this.rtcSession.dataChannel.send(JSON.stringify(pingMsg));
        }
    },

    openSettingDrawer() {
      this.settingDrawerOpen = true;
      this.listKeyboardBindings();
    },

    onDataChannelMsg(msg) {
      const msgStr = String.fromCharCode.apply(null, new Uint8Array(msg.data));
      const msgObj = JSON.parse(msgStr);
      switch (msgObj.type) {
        case MessagePing:
          this.stats.rtt = new Date().getTime() - msgObj.timestamp;
          break;
        case MessageChat:
          if (!msgObj["from"]) return;
          api.get("api/v1/user/"+msgObj["from"]).then(resp=>{
            notification.info({
              message: resp["data"]["name"],
              description: msgObj.data,
              placement: "topLeft",
              duration: 1,
            });
          });
          break;
        default:
          break
      }
    },

    onDataChannelOpen() {
      this.chatBtnDisabled = false;
      const _this = this;
      this.pingInterval = setInterval(_=>{
        _this.ping();
        _this.collectRTCStats();
      }, 1000);
    },

    onDataChannelClose() {
        if (this.pingInterval) clearInterval(this.pingInterval);
        this.chatBtnDisabled = true;
    },

    listKeyboardBindings: function() {
      const _this = this;
      api.get("api/v1/keyboard/bindings?page=0&pageSize=100").then(resp=>{
        _this.configs.userBindings = resp["bindings"];
        _this.configs.userBindings.push(defaultBinding);
        const userBindings = _this.configs.userBindings;
        let options = [];
        for (let i = 0; i < userBindings.length; i++) {
          options.push({
            value: userBindings[i]["id"],
            label: userBindings[i]["name"],
          });
          const curBindings = userBindings[i]["bindings"];
          for (let j = 0; j < curBindings.length; j++) {
            curBindings[j]["buttons"] = [curBindings[j]["keyboardKey"]];
          }
        }
        _this.configs.selectedBindingKey = userBindings[0].id;
        _this.configs.selectedBinding = userBindings[0];
        _this.configs.bindingOptions = options;
      }).catch(_=>{
        message.error("获取按键绑定失败");
      });
    },
    onKeyboardBindingSelectChange: function(ev) {
      this.configs.selectedBinding = this.configs.userBindings.find(item=>item["id"]===this.configs.selectedBindingKey);
    },

    updateBinding: function() {
      const data = this.convertBindingsToApiObj(this.configs.selectedBinding);
      data["id"] = this.configs.selectedBindingKey;
      const _this = this;
      this.createBtnDisable = true;
      this.bindingBtnDisable = true;
      api.put("api/v1/keyboard/binding", data).then(_=>{
        message.success("修改成功");
        _this.listKeyboardBindings();
        _this.createBtnDisable = false;
        _this.bindingBtnDisable = false;
      }).catch(_=>{
        message.error("修改失败");
        _this.createBtnDisable = false;
        _this.bindingBtnDisable = false;
      });
    },

    convertBindingsToApiObj: function(bindingObj) {
      const apiObj = {
        "name": bindingObj.name,
        "bindings": [],
      };
      bindingObj.bindings.forEach(item=>{
        apiObj["bindings"].push({
          "emulatorKey": item.emulatorKey,
          "keyboardKey": item.buttons[0],
        });
      });
      return apiObj;
    },

    setKeyboardBindingEnabled: function() {
      this.setKeyboardControl(true);
      message.info("设置成功");
    },

    collectRTCStats: function() {
      const _this = this;
      if (this.rtcSession && this.rtcSession.pc) {
        const pc = this.rtcSession.pc;
        pc.getStats().then(stats=>{
          stats.forEach(report=>{
            if (report.type === "inbound-rtp" && report.kind === "video") {
              _this.stats.fps = report.framesPerSecond;
            }
          });
        });
      }
    },
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
#stats {
  position: absolute;
  right: 0;
  top: 0;
}
</style>