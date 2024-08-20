<template>
    <a-row>
        <a-col :xs="{ offset: 2, span: 20 }" :sm="{ offset: 2, span: 20 }" :md="{ offset: 6, span: 12 }"
            :lg="{ offset: 8, span: 8 }">
            <a-card :bordered="false">
                <a-row>
                    <a-col :offset="4" :span="16">
                        <h1>NESGO</h1>
                    </a-col>
                </a-row>
                <a-form layout="vertical" :model="formState" name="basic" :label-col="{ span: 4 }" autocomplete="off"
                    @finish="onFinish" @finishFailed="onFinishFailed">
                    <a-form-item label="用户名" name="name" :rules="[{ required: true, message: '请输入用户名' }]">
                        <a-input v-model:value="formState.name" />
                    </a-form-item>

                    <a-form-item label="密码" name="password" :rules="[{ required: true, message: '请输入密码' }]">
                        <a-input-password v-model:value="formState.password" />
                    </a-form-item>
                    <a-row>
                        <a-col :offset="4" :span="16" style="text-align: center">
                            新用户？点击<RouterLink to="/register">此处</RouterLink>注册
                        </a-col>
                    </a-row>
                    <a-row>
                        <a-col :offset="4" :span="16">
                            <a-button type="primary" id="loginButton" html-type="submit">登录</a-button>
                        </a-col>
                    </a-row>
                </a-form>
            </a-card>

        </a-col>
    </a-row>
</template>

<script>
import { Col, Row, Card } from 'ant-design-vue';
import { Button, Form, Input } from 'ant-design-vue';
import api from '../api/request';
import tokenStorage from '../api/token';
import { message } from 'ant-design-vue';
import router from '../router';
import { RouterLink } from "vue-router";

export default {
    components: {
        ACard: Card,
        ARow: Row,
        ACol: Col,
        AButton: Button,
        AForm: Form,
        AFormItem: Form.Item,
        AInput: Input,
        AInputPassword: Input.Password,
        RouterLink,
    },
    data() {
        return {
            formState: {
                name: "",
                password: "",
            }
        }
    },
    methods: {
        onFinish(ev) {
            api.post("api/v1/login", this.formState)
                .then(resp => {
                    message.success("登录成功")
                    tokenStorage.setToken(resp["token"])
                    router.push("/home")
                }).catch(_ => {
                    message.warn("登录失败，请检查用户名和密码");
                })
        },
        onFinishFailed(ev) {

        },
    }
}
</script>

<style>
body {
    background-color: #CCCCCC;
}

#loginButton {
    width: 100%;
}
</style>