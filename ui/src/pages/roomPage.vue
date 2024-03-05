<template>
  <a-row style="height: 100vh; background-color: grey">
    <a-col :span="8">
      <a-button @click="_=>{membersDrawerOpen=true}"></a-button>
    </a-col>
    <a-col :span="8" style="height: 100vh">
      <a-card style="height: 100%">
        <video width="100%" id="video"></video>
        <a-row>
          <a-col :span="8">
            <a-button style="width: 90%;" type="primary" @click="connect" :disabled="connectBtnDisabled">连接</a-button>
          </a-col>
          <a-col :span="8" :offset="8">
            <a-button style="width: 90%;" type="primary">重启</a-button>
          </a-col>
        </a-row>
        <a-row style="margin-top: 20px">
          <a-col :span="8">
            <a-button style="width: 90%;" type="primary">保存</a-button>
          </a-col>
          <a-col :span="8" :offset="8">
            <a-button style="width: 90%;" type="primary">加载</a-button>
          </a-col>
        </a-row>
      </a-card>
    </a-col>
    <a-col :span="8"></a-col>
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
                <a-radio-group v-model:value="item.role" :disabled="memberSelf.role!==0" >
                  <a-radio :value="1">玩家</a-radio>
                  <a-radio :value="2">观战</a-radio>
                </a-radio-group>
              </a-col>
              <a-col :span="6">
                <a-button type="primary" :hidden="memberSelf.role!==0">踢出</a-button>
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
      window.onorientationchange = _=>{
        if (window.orientation === 180 || window.orientation === 0) {
          message.info("手机横屏以获取最佳体验")
        }
      }
      this.roomId = this.$route["params"]["roomId"]
      this.listRoomMembers()
      this.getMemberSelf()
  },
  methods: {
      listRoomMembers() {
        api.get("/room/" + this.roomId + "/members")
            .then(resp=>{
              this.members = resp.data
              this.members.forEach(m=>{
                m.controls = this.p1p2Ratio(m)
              })
            })
      },
      getMemberSelf() {
        api.get("/room/" + this.roomId + "/member")
            .then(resp=>{
              this.memberSelf = resp.data
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
              .then(_=>this.onConnected())
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
              urls: this.rtcSession.turnAddress,
              username: this.rtcSession.turnUser,
              credential: this.rtcSession.turnPassword,
            },
          ],
          iceTransportPolicy: "relay",
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
            this.rtcSession.ws.close()
            break
          case "disconnected":
            pc.close()
            this.connectBtnDisabled = false
            break
          default:
            break
        }
      },
      onConnected() {
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
        }
        // todo screen buttons
        // todo enable restart save load buttons
      },
    sendAction(code, pressed) {
      this.rtcSession.dataChannel.send(JSON.stringify({
        "type": pressed,
        "data": btoa(code),
      }))
    }
  }
}
</script>

<style scoped>
#video {
  position: inherit;
  margin: auto;
  top: 0;
  bottom: 0;
  background-color: black;
}
</style>