<template>
    <a-layout style="height: 100vh;">
        <a-layout-header style="position: fixed; z-index: 1; width: 100%;">
            <a-menu theme="dark" mode="horizontal" v-model:selectedKeys="headerSelectedKeys">
                <a-menu-item key="1">我的房间</a-menu-item>
                <a-menu-item key="2">加入房间</a-menu-item>
                <a-menu-item key="3">设置</a-menu-item>
            </a-menu>
        </a-layout-header>
        <a-row style="margin-top: 64px;">
            <a-col :xs="{offset: 2, span: 20}" :sm="{offset: 2, span: 20}" :md="{offset: 4, span: 16}" :lg="{offset: 4, span: 16}">
                <RoomList v-if="headerSelectedKeys[0] === '1'" :joined="true"/>
                <RoomList v-else-if="headerSelectedKeys[0] === '2'" :joined="false"></RoomList>
                <UserSetting v-else-if="headerSelectedKeys[0] === '3'"></UserSetting>
            </a-col>
        </a-row>
    </a-layout>
</template>

<script>
import { Layout, Menu, Card } from 'ant-design-vue';
import { Row, Col } from "ant-design-vue";
import RoomList from "../components/roomList.vue";
import UserSetting from "../components/userSetting.vue";
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
        RoomList: RoomList,
        UserSetting: UserSetting,
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