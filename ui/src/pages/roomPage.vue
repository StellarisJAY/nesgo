<template>
  <a-row style="height: 100vh; background-color: grey">
    <!--left side buttons-->
    <a-col :span="6">
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
    <a-col :span="12" style="height: 100vh">
      <a-card style="height: 100%">
        <a-row style="height: 80%">
          <video id="video" playsinline webkit-playsinline="true"></video>
        </a-row>
        <a-row>
          <a-col :span="8">
            <a-button style="width: 90%;" type="primary" @click="connect" :disabled="connectBtnDisabled">连接</a-button>
          </a-col>
          <a-col :span="8" :offset="8">
            <a-button style="width: 90%;" type="primary" :disabled="restartBtnDisabled">重启</a-button>
          </a-col>
        </a-row>
        <a-row style="margin-top: 20px">
          <a-col :span="8">
            <a-button style="width: 90%;" type="primary" :disabled="saveBtnDisabled">保存</a-button>
          </a-col>
          <a-col :span="8" :offset="8">
            <a-button style="width: 90%;" type="primary" :disabled="loadBtnDisabled">加载</a-button>
          </a-col>
        </a-row>
      </a-card>
    </a-col>
    <!--right side buttons-->
    <a-col :span="6">
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
    <a-drawer
        v-model:open="membersDrawerOpen"
        :root-style="{ color: 'blue' }"
        title="成员列表"
        placement="right"
    >
      <a-list item-layout="vertical" :data-source="members">
        <template #renderItem="{item}">
          <a-list-item>
            <a-row>
              <a-col :span="8"><CrownTwoTone v-if="item.role===0" />{{item.name}}</a-col>
              <a-col :span="10">
                <a-radio-group v-model:value="item.controls" :disabled="memberSelf.role!==0">
                  <a-radio value="1">P1</a-radio>
                  <a-radio value="2">P2</a-radio>
                </a-radio-group>
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
  </a-row>
</template>

<script>
import api from "../api/request.js";
import { Row, Col } from "ant-design-vue";
import {CrownTwoTone} from '@ant-design/icons-vue';
import {Card, Button, Drawer, List, Input, RadioGroup, Radio} from "ant-design-vue";
import {message} from "ant-design-vue";
import tokenStorage from "../api/token.js";
import router from "../router/index.js";
import {ArrowUpOutlined, ArrowDownOutlined, ArrowLeftOutlined, ArrowRightOutlined} from "@ant-design/icons-vue"

const MessageSDPOffer = 0
const MessageSDPAnswer = 1
const MessageICECandidate = 2
const MessageGameButtonPressed = 3
const MessageGameButtonReleased = 4
const MessageTurnServerInfo = 5

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
              existingGames: {},
            }

        }
    },
  created() {
      this.roomId = this.$route["params"]["roomId"]
      this.getMemberSelf()
  },
  beforeDestroy() {
    this.rtcSession.ws.close()
    this.rtcSession.pc.close()
  },
  methods: {
      openRoomMemberDrawer() {
          this.listRoomMembers().then(_=>this.membersDrawerOpen=true)
      },
      listRoomMembers() {
        return api.get("/room/" + this.roomId + "/members")
            .then(resp=>{
              this.members = resp.data
              this.members.forEach(m=>{
                m.controls = this.p1p2Ratio(m)
              })
            })
            .catch(resp=>{
              if (resp.status === 403) {
                message.warn("你不是该房间成员")
                router.push("/home")
              }
            })
      },
      getMemberSelf() {
        api.get("/room/" + this.roomId + "/member")
            .then(resp=>{
              this.memberSelf = resp.data
            }).catch(resp=>{
              if (resp.status === 403) {
                message.warn("你不是该房间成员")
                router.push("/home")
              }
        })
      },
      p1p2Ratio(m) {
        if (m["player1"]) {
          return 1
        }else if (m["player2"]) {
          return 2
        }
        return 0
      },

      connect() {
        this.connectBtnDisabled = true
        const ws = api.webSocket("/room/"+this.roomId+"/rtc?auth=" + tokenStorage.getToken() + "&game="+this.selectedGame)
        ws.onclose = ev => {
          console.log("websocket connection closed")
        }
        ws.onerror = ev => {
          ws.close()
        }
        ws.onmessage = this.onWebsocketMessage
        this.rtcSession.ws = ws
      },

      onWebsocketMessage(ev) {
        const message = JSON.parse(ev.data)
        const ws = this.rtcSession.ws
        if (message["type"] === MessageSDPOffer) {
          const sdpOffer = JSON.parse(atob(message["data"]))
          const pc = this.rtcSession.pc
          pc.setRemoteDescription(sdpOffer)
              .then(_ => pc.createAnswer())
              .then(sdp => pc.setLocalDescription(sdp))
              .then(_ => {
                console.log("remote sdp: ", pc.remoteDescription)
                console.log("local sdl:  ", pc.localDescription)
                ws.send(JSON.stringify({
                  "type": MessageSDPAnswer,
                  "data": btoa(JSON.stringify(pc.localDescription)),
                }))
                pc.addTransceiver("video")
              })
              .catch(err=>{
                console.log(err)
              })
        }else if (message["type"] === MessageTurnServerInfo) {
          const turnInfo = JSON.parse(atob(message["data"]))
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
          datachannel.onclose = _=>{
            window.onkeydown = _=>{}
            window.onkeyup = _ => {}
          }
          datachannel.onerror = err=>{
            console.log(err)
          }
          datachannel.onmessage = msg=>{
            console.log("unexpected dataChannel message:", msg)
          }
        }
        this.rtcSession.pc = pc
      },
      onICECandidate(ev) {
        if (ev.candidate !== null) {
          this.rtcSession.ws.send(JSON.stringify({
            "type": MessageICECandidate,
            "data": btoa(JSON.stringify(ev.candidate))
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
        message.success("连接成功")
        if (this.memberSelf["role"] !== RoleObserver) {
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
          // todo screen buttons
          // todo enable restart save load buttons
          this.initControlButtons()
        }
      },
    onDisconnected() {
      this.disableControlButtons()
      this.rtcSession.pc.close()
      this.connectBtnDisabled = false
    },
    sendAction(code, pressed) {
      this.rtcSession.dataChannel.send(JSON.stringify({
        "type": pressed,
        "data": btoa(code),
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
</style>