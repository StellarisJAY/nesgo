import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: "/login",
            name: "login",
            component: ()=>import("../pages/loginPage.vue")
        },
        {
            path: "/register",
            name: "register",
            component: ()=>import("../pages/registerPage.vue")
        }
    ]
})

export default router;