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
            <a-button ref="refConnBtn" style="width: 90%;" type="primary" @click="connect" :disabled="connectBtnDisabled">
              <a-spin v-if="loading" :spinning="loading"></a-spin>
              <p v-else>连接</p>
            </a-button>
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
      <a-form v-if="memberSelf.role === 0">
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
              <a-col :span="8"><CrownTwoTone v-if="item.role===0" />{{item.name}}</a-col>
              <a-col :span="10">
                <a-checkbox :disabled="memberSelf.role!==0" v-model:checked="item['player1']" @change="ev=>{onP1P2Change(ev, item, 1)}">P1</a-checkbox>
                <a-checkbox :disabled="memberSelf.role!==0" v-model:checked="item['player2']" @change="ev=>{onP1P2Change(ev, item, 2)}">P2</a-checkbox>
                <a-radio-group v-model:value="item.role" :disabled="memberSelf.role!==0 || item.role===0"
                               @change="ev=>{onRoleRatioChange(ev, item)}">
                  <a-radio :value="1">玩家</a-radio>
                  <a-radio :value="2">观战</a-radio>
                </a-radio-group>
              </a-col>
              <a-col :span="6">
                <a-button type="primary" :hidden="memberSelf.role!==0"
                          :disabled="item.role===0 || memberSelf.role!==0"
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
            <a-button type="primary" @click="loadSavedGame(item.id)" :disabled="memberSelf.role!==0">加载</a-button>
            <a-button danger @click="deleteSavedGame(item.id)" :disabled="memberSelf.role!==0">删除</a-button>
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
import {notification, Tour, Spin} from "ant-design-vue";

const MessageSDPOffer = 0
const MessageSDPAnswer = 1
const MessageICECandidate = 2
const MessageGameButtonPressed = 3
const MessageGameButtonReleased = 4
const MessageTurnServerInfo = 5
const MessageChat = 6
const MessagePing = 7
const MessagePong = 8

const RoleHost = 0
const RoleGamer = 1
const RoleObserver = 2

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
    ASpin: Spin,
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
            selectedGame: "SuperMario.nes",
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
              existingGames: [{
                value: "SuperMario.nes",
                label: "SuperMario.nes",
              }],
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
            loading: false,
        }
    },
  created() {
      this.roomId = this.$route["params"]["roomId"]
      this.getMemberSelf().then(_=>{
        if (this.memberSelf.role === RoleHost) {
          this.getRoomInfoWithPassword()
        }
      })
      this.listRoomMembers()
      this.listGames()
  },
  unmounted() {
    this.rtcSession.ws.close()
    this.rtcSession.pc.close()
  },
  methods: {
      openRoomMemberDrawer() {
          this.listRoomMembers().then(_=>this.membersDrawerOpen=true)
      },
      getRoomInfoWithPassword() {
        api.get("/room/" + this.roomId + "/fullInfo").then(resp=>{
          if (resp.status === 200) {
            this.fullRoomInfo = resp.data
          }
        })
      },
      listRoomMembers() {
        return api.get("/room/" + this.roomId + "/members")
            .then(resp=>{
              this.members = resp.data
            })
            .catch(resp=>{
              if (resp.status === 403) {
                message.warn("你不是该房间成员")
                router.push("/home")
              }
            })
      },
      getMemberSelf() {
        return api.get("/room/" + this.roomId + "/member")
            .then(resp=>{
              this.memberSelf = resp.data
            }).catch(resp=>{
              if (resp.status === 403) {
                message.warn("你不是该房间成员")
                router.push("/home")
              }
        })
      },

      connect() {
        this.loading = true
        this.connectBtnDisabled = true
        const ws = api.webSocket("/room/"+this.roomId+"/rtc?auth=" + tokenStorage.getToken() + "&game="+this.selectedGame)
        ws.onclose = ev => {
          console.log("websocket connection closed")
        }
        ws.onerror = ev => {
          message.error("连接出错")
          ws.close()
          if (this.rtcSession.pc) {
            this.rtcSession.pc.close()
          }
          router.push("/home")
        }
        ws.onmessage = this.onWebsocketMessage
        this.rtcSession.ws = ws
      },

      onWebsocketMessage(ev) {
        const message = JSON.parse(ev.data)
        const ws = this.rtcSession.ws
        if (message["type"] === MessageSDPOffer) {
          const sdpOffer = JSON.parse(message["data"])
          const pc = this.rtcSession.pc
          pc.setRemoteDescription(sdpOffer)
              .then(_ => pc.createAnswer())
              .then(sdp => pc.setLocalDescription(sdp))
              .then(_ => {
                console.log("remote sdp: ", pc.remoteDescription)
                console.log("local sdl:  ", pc.localDescription)
                ws.send(JSON.stringify({
                  "type": MessageSDPAnswer,
                  "data":JSON.stringify(pc.localDescription),
                }))
                pc.addTransceiver("video")
              })
              .catch(err=>{
                console.log(err)
              })
        }else if (message["type"] === MessageTurnServerInfo) {
          const turnInfo = JSON.parse(message["data"])
          this.rtcSession.turnAddress = turnInfo["address"]
          this.rtcSession.turnUser = turnInfo["username"]
          this.rtcSession.turnPassword = turnInfo["password"]
          this.createPeerConnection()
        }
      },

      createPeerConnection() {
        const ws = this.rtcSession.ws
        const pc = new RTCPeerConnection({
          iceServers: [
            {
              urls: "stun:192.168.0.107:3478"
            },
            {
              urls: this.rtcSession.turnAddress,
              username: this.rtcSession.turnUser,
              credential: this.rtcSession.turnPassword,
            },
          ],
          iceTransportPolicy: "relay"
        })

        pc.onicecandidate = this.onICECandidate
        pc.onconnectionstatechange = this.onPeerConnStateChange
        pc.oniceconnectionstatechange = ev=>{
          console.log("ice conn state: " + pc.iceConnectionState)
        }

        pc.ontrack = ev=>{
          document.getElementById("video").srcObject = ev.streams[0]
          document.getElementById("video").autoplay = true
          document.getElementById("video").controls = true
        }

        pc.ondatachannel = ev=>{
          const datachannel = ev.channel
          this.rtcSession.dataChannel = datachannel
          this.chatBtnDisabled = false
          datachannel.onclose = _=>{
            clearInterval(this.pingInterval)
            this.chatBtnDisabled = true
            window.onkeydown = _=>{}
            window.onkeyup = _ => {}
          }
          datachannel.onerror = err=>{
            console.log(err)
          }
          datachannel.onmessage = msg=>{
            const message = JSON.parse(msg.data)
            if (message.type === MessageChat) {
              this.onChatMessage(JSON.parse(message.data))
            }else if (message.type === MessagePong) {
              const now = new Date().getTime()
              this.rtt  = now - Number(message.data)
            }
          }
          this.pingInterval = setInterval(_=>{this.ping()}, 3000);
        }
        this.rtcSession.pc = pc
      },
      onICECandidate(ev) {
        if (ev.candidate !== null) {
          this.rtcSession.ws.send(JSON.stringify({
            "type": MessageICECandidate,
            "data": JSON.stringify(ev.candidate)
          }))
          this.rtcSession.pc
              .addIceCandidate(ev.candidate)
              .then(_ => console.log("ice candidate: ", ev.candidate))
        }
      },
      onPeerConnStateChange(_) {
        const pc = this.rtcSession.pc
        console.log("peer conn state: " + pc.connectionState)
        switch (pc.connectionState) {
          case "connected":
            this.onConnected()
            this.rtcSession.ws.close()
            break
          case "disconnected":
            this.onDisconnected()
            break
          default:
            break
        }
      },
      onConnected() {
        this.loading = false
        message.success("连接成功")
        this.setKeyboardControl(true)
        this.initControlButtons()
        this.saveBtnDisabled = this.memberSelf["role"] !== RoleHost
        this.loadBtnDisabled = false
        this.restartBtnDisabled = this.memberSelf["role"] !== RoleHost
      },
    onDisconnected() {
      this.loading = false
      this.disableControlButtons()
      this.rtcSession.pc.close()
      this.connectBtnDisabled = false
    },
    sendAction(code, pressed) {
      this.rtcSession.dataChannel.send(JSON.stringify({
        "type": pressed,
        "data": code,
      }))
    },

    onRoleRatioChange(ev, member) {
        member.role = ev.target.value
        api.post("/room/" + this.roomId + "/role", {
            "memberId": member.id,
            "role": member.role,
            "kick": false,
        })
            .then(resp=>{
                if (resp.status === 200) {
                  message.success("修改成功")
                }
            })
    },
    kickMember(member) {
        api.post("/room/" + this.roomId + "/member/kick", {
          "memberId": member.id,
          "kick": true,
        }).then(resp=>{
            if (resp.status === 200) {
              message.success("成功")
              this.members = this.members.filter(item=>{
                return item.id !== member.id
              })
            }
        })
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
      return api.get("/room/" + this.roomId + "/saves").then(resp=>{
          if (resp.status === 200) {
            this.savedGames = resp.data
          }
      })
    },
    openSavedGamesDrawer() {
        this.getSavedGames().then(_=>this.savedGameOpen=true)
    },
    saveGame() {
        api.post("/room/" + this.roomId + "/quickSave").then(resp=>{
          if (resp.status === 200) {
            message.success("保存成功")
          }
        }).catch(resp=>{
          if (resp.status === 400) {
            message.error("无法保存：已达到存档数量上限")
          }
        })
    },
    loadSavedGame(id) {
        api.post("/room/" + this.roomId + "/load/" + id).then(resp=>{
          if (resp.status === 200) {
            message.success("加载存档成功")
          }
        }).catch(resp=>{
          message.error(resp.message)
        })
    },
    deleteSavedGame(id) {
        api.post("/room/" + this.roomId + "/saves/" + id + "/delete").then(resp=>{
          if (resp.status === 200) {
            message.success("删除成功")
            this.savedGames = this.savedGames.filter(s=> {
              return s.id !== id
            })
          }
        }).catch(resp=>{
          message.error(resp.message)
        })
    },
    restart() {
        const game = this.selectedGame
        api.post("/room/" + this.roomId + "/restart?game=" + game).then(resp=>{
          if (resp.status === 200) {
            message.success("重启模拟器成功")
          }
        })
    },
    listGames() {
      api.get("/games").then(resp=>{
        resp.data.forEach(game=>{
          this.configs.existingGames.push({
            value: game.name,
            label: game.name,
            data: game,
          })
        })
      })
    },
    onP1P2Change(ev, m, which) {
      if (which === 1) {
        if (ev.target.checked && m["player2"]) {
          m["player2"] = false
        }
      }else {
        if (ev.target.checked && m["player1"]) {
          m["player1"] = false
        }
      }
      api.post("/room/" + this.roomId + "/control/transfer", {
        "memberId": m.id, "setController1": m["player1"], "setController2": m["player2"]}).then(resp=>{
          if (resp.status ===  200) {
            message.success("修改成功")
          }
      }).catch(resp=>{
          if (resp.status === 400 || resp.status === 404) {
            message.error("无法转移控制权：用户未连接")
            m.player1 = false
            m.player2 = false
          }else if (resp.status === 403) {
            message.error("无法转移控制权：用户为观战者")
            m.player1 = false
            m.player2 = false
          }
      })
    },
    alterRoomPrivacy() {
      api.post("/room/" + this.roomId + "/alter", {"name": this.fullRoomInfo.name, "private": this.fullRoomInfo.private}).then(resp=>{
        if (resp.status === 200) {
          message.success("修改成功")
          this.fullRoomInfo = resp.data
        }
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
      this.rtcSession.dataChannel.send(JSON.stringify({
        "type": MessageChat,
        "data": this.chatMessage}))
    },
    onChatMessage(msg) {
        const m = this.members.find(m=>{return m.id === msg.from})
        if (!m) {
          this.listRoomMembers()
        }
        notification.open({
          message: m ? m.name : "ID:"+m.id,
          description: msg.content,
          duration: 1,
          placement: "topRight",
        })
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
        this.rtcSession.dataChannel.send(JSON.stringify({
          "type": MessagePing,
          "data": String(timestamp),
        }))
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