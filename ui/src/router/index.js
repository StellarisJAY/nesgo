import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: "/",
            name: "root",
            component: ()=>import("../pages/homePage.vue")
        },
        {
            path: "/login",
            name: "login",
            component: ()=>import("../pages/loginPage.vue")
        },
        {
            path: "/register",
            name: "register",
            component: ()=>import("../pages/registerPage.vue")
        },
        {
            path: "/home",
            name: "home",
            component: ()=>import("../pages/homePage.vue")
        },
        {
            path: "/room/:roomId",
            name: "room",
            component: ()=>import("../pages/roomPage.vue")
        }
    ]
})

export default router;