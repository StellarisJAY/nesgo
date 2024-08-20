<template>
    <a-row>
        <a-col :xs="{ offset: 2, span: 20 }" :sm="{ offset: 2, span: 20 }" :md="{ offset: 6, span: 12 }"
            :lg="{ offset: 8, span: 8 }">
            <a-card :bordered="false">
                <a-row>
                    <a-col :offset="4" :span="16">
                        <h1>新用户注册</h1>
                    </a-col>
                </a-row>
                <a-form layout="vertical" :model="formState" name="basic" :label-col="{ span: 4 }" autocomplete="off"
                    @finish="onFinish" @finishFailed="onFinishFailed">
                    <a-form-item label="用户名" name="name" :rules="[{ required: true, message: '请输入用户名' }]">
                        <a-input v-model:value="formState.name" />
                    </a-form-item>

                    <a-form-item label="密码" name="password" :rules="rules.pass">
                        <a-input-password v-model:value="formState.password" />
                    </a-form-item>

                    <a-form-item label="确认密码" name="confirm" :rules="rules.confirmPass">
                        <a-input-password v-model:value="formState.confirmPassword" />
                    </a-form-item>
                    <a-row>
                        <a-col :span="16" :offset="4" style="text-align: center">
                            已经拥有账号？点击<RouterLink to="/login">此处</RouterLink>登录
                        </a-col>
                    </a-row>
                    <a-row>
                        <a-col :offset="4" :span="16">
                            <a-button type="primary" style="width: 100%;" html-type="submit">注册</a-button>
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
                confirmPassword: "",
            },
            rules: {
                pass: [
                    {
                        required: true,
                        validator: this.validatePass,
                        trigger: 'change',
                    },
                ],
                confirmPass: [
                    {
                        required: true,
                        validator: this.validatePass2,
                        trigger: 'change',
                    },
                ],
            }
        }
    },
    methods: {
        onFinish(ev) {
            api.post("api/v1/register", {
                "name": this.formState.name,
                "password": this.formState.password,
            })
                .then(data => {
                    message.success("注册成功")
                    router.push("/login")
                })
                .catch(resp => {
                })
        },
        onFinishFailed(ev) {

        },
        validatePass() {
            if (this.formState.password === '') {
                return Promise.reject("请输入密码")
            }
            return Promise.resolve()
        },
        validatePass2() {
            if (this.formState.confirmPassword === '') {
                return Promise.reject("请确认密码")
            } else if (this.formState.confirmPassword !== this.formState.password) {
                return Promise.reject("两次输入密码不同")
            }
            return Promise.resolve()
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