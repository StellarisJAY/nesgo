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
        },
        {
            path: "/home",
            name: "home",
            component: ()=>import("../pages/homePage.vue")
        }
    ]
})

export default router;