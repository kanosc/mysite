import { createRouter, createWebHistory } from 'vue-router'

import Home from '../views/Home.vue'

const m = new Map([
    ['username', 'please input your username'],
    ['password', 'please input password of the room']
])

const routes = [
    { path: '/', name: 'Home', component: Home },
    { path: '/filesys', name: 'FileSys', component: () => import('../views/FileSys.vue') },
    // {
    //     path: '/login/:option',
    //     name: 'Login',
    //     component: () => import('../views/Login.vue'),
    //     props: route => ({ ...route.params, option: route.params.option, info: m.get(route.params.option) })
    // },
    {
        path: '/login/username',
        name: 'LoginUser',
        component: () => import('../views/Login.vue')
    },
    {
        path: '/login/password/:roomname',
        name: 'LoginPassword',
        component: () => import('../views/Password.vue'),
        props: route => ({ ...route.params, roomname: route.params.roomname })
    },
    {
        path: '/page/room/:roomname',
        name: 'Room',
        component: () => import('../views/ChatRoom.vue'),
        props: route => ({ ...route.params, roomname: route.params.roomname })
    },
    {
        path: '/chatrooms',
        name: 'chatrooms',
        component: () => import('../views/Chat.vue'),
        meta: {
            requiresUsername: true
        }
    },
    {
        path: '/keysignup',
        name: 'keysignup',
        component: () => import('../views/KeySignUp.vue'),
    },
    { path: '/fortune', name: 'Fortune', component: () => import('../views/FortuneCard.vue') },
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('../views/NotFound.vue')
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach((to, from) => {
    let un = localStorage.getItem("username-chat")
    if (to.meta.requiresUsername && un == null) {
        return { path: '/login/username' }
    }
})

export default router