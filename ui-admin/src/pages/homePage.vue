<template>
    <a-layout style="height: 100vh;">
        <a-layout-header style="position: fixed; z-index: 1; width: 100%;">
            <a-menu theme="dark" mode="horizontal" v-model:selectedKeys="headerSelectedKeys">
                <a-menu-item key="1">模拟器</a-menu-item>
                <a-menu-item key="2">房间管理</a-menu-item>
                <a-menu-item key="3">游戏管理</a-menu-item>
                <a-menu-item key="4">管理员</a-menu-item>
            </a-menu>
        </a-layout-header>
        <a-row style="margin-top: 64px;">
            <a-col :xs="{offset: 2, span: 20}" :sm="{offset: 2, span: 20}" :md="{offset: 4, span: 16}" :lg="{offset: 4, span: 16}">
                <EmulatorList v-if="headerSelectedKeys[0]==='1'"></EmulatorList>
                <GameList v-if="headerSelectedKeys[0]==='3'"></GameList>
            </a-col>
        </a-row>
    </a-layout>
</template>

<script>
import { Layout, Menu, Card, Button } from 'ant-design-vue';
import { Row, Col } from "ant-design-vue";
import EmulatorList from "../components/emulatorList.vue"
import GameList from "../components/gameList.vue"

export default {
    components: {
        ALayout: Layout,
        ALayoutHeader: Layout.Header,
        AMenu: Menu,
        AMenuItem: Menu.Item,
        ALayoutContent: Layout.Content,
        ARow: Row,
        ACol: Col,
        ACard: Card,
        AButton: Button,
        EmulatorList: EmulatorList,
        GameList: GameList,
    },
    data() {
        return {
            headerSelectedKeys: ['1']
        }
    }
}
</script>

<style>
</style>