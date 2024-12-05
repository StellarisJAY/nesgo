## NESGO
<a href="#" target="_blank">
    <img alt="Go version" src="https://badgen.net/static/Go/1.21"/>
</a>
<a href="https://goreportcard.com/report/github.com/stellarisjay/nesgo" target="_blank">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/stellarisjay/nesgo" />
</a>
<a href="#" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
</a>
<a href="https://webrtc.github.io/webrtc-org/start/" target="_blank">
    <img alt="webrtc" src="https://img.shields.io/badge/WebRTC-333?logo=webrtc&logoColor=0da&style=static"/>
</a>
<a href="https://github.com/pion/webrtc" target="_blank">
    <img alt="pion/webrtc" src="https://badgen.net/static/pion/webrtc?color=yellow"/>
</a>

**停止更新，请前往新项目：https://github.com/StellarisJAY/cloud-emu**

基于WebRTC的云游戏NES模拟器，支持最多4人同屏联机游玩（2玩家+2观战）。


### 功能列表

- [x] 多人游戏房间，房主可设置权限。
- [x] 模拟器热重启，重启模拟器不需要断开游戏连接。
- [x] 保存与读取，跨游戏加载存档自动重启模拟器
- [x] 房间内即时聊天
- [x] 自定义按键设置
- [x] 管理员-游戏上传
- [x] 模拟器服务水平扩展
- [x] 画面设置
- [x] 游戏加速
- [ ] 存档转移（跨房间转移、上传下载）
- [ ] 更多模拟器

## Issues

- 暂不支持部分卡带格式
- 部分游戏存在贴图渲染错误
- 部分卡带格式暂不支持存档

## 安装部署

### 编译运行

安装依赖：pkg-config libx264 libvpx libopus

libx264和libvpx提供视频编码，libopus提供音频编码

```shell
apt install pkg-config libx264-dev libopusfile-dev libvpx-dev
```

编译后端服务

```shell
# 编译微服务后端
cd backend
go mod tidy
make build # 可执行文件将在bin目录下
```
创建配置文件，请参考/backend/configs目录

运行后端服务
```shell
./bin/gaming --conf ./configs/gaming.yaml
./bin/user   --conf ./configs/user.yaml
./bin/room   --conf ./configs/room.yaml
./bin/webapi --conf ./configs/webapi.yaml
./bin/admin  --conf ./configs/admin.yaml
```
编译运行前端
```shell
cd ui 
yarn build
# 获得dist目录，可使用nginx部署
```
