<script setup>
import { ref } from 'vue'
import { copyList, convertSelectToColor } from '../common/common'
import { useRouter } from 'vue-router'
const router = useRouter()
const props = defineProps(['roomname'])
const roomname = props.roomname
const username = localStorage.getItem('username-chat')

const inputText = ref('')
const messageData = ref('')
let originMessages = []

let userColorTheme = localStorage.getItem('color-theme')
console.log(`color theme: ${userColorTheme}`)
if (userColorTheme == null) {
    userColorTheme = 'Dark'
}

const getLocalToken = (roomname) => {
    return localStorage.getItem(roomname + '-token')
}

const getMessages = async () => {
    try {
        let localToken = getLocalToken(roomname)
        let fetchOpt = {}
        if (localToken != null) {
            fetchOpt.headers = {
                'Authorization': localToken
            }
        }
        const res = await fetch("/api/chat/v1/msg/" + roomname, fetchOpt)
        if (res.status == 401) {
            router.push("/login/password/" + roomname)
            return
        }
        const rsp = await res.json()
        originMessages = rsp.messages
        messageData.value = copyList(originMessages).map(parseMessage)
    } catch (e) {
        alert(e)
    }
}
getMessages()

const getWebSocketProtocol = () => {
    if (window.location.protocol == 'https:') {
        return 'wss'
    }
    else {
        return 'ws'
    }
}
const wsurl = getWebSocketProtocol() + "://" + window.location.host + "/api/chat/v1/wschat/" + roomname
console.log(wsurl)
console.log(window.location.protocol)
let ws = null
const createWebSocket = (url) => {
    try {
        console.log("create socket begin");

        if (window.ws == null) {
            ws = new WebSocket(url)
        } else {
            console.log('websocket object already exist')
            console.log(`status of websocket is ${ws.readyState}`)
        }
        console.log("create socket end");
    } catch (e) {
        alert(e)
    }
}

const setWebSocketCallback = (ws) => {
    ws.onmessage = (msg) => {
        originMessages.unshift(parseMessage(msg.data))
        messageData.value.unshift(parseMessage(msg.data))
        console.log(`receive msg ${msg.data}`)
    }
    ws.onopen = () => {
        console.log('websocket opened')
    }
    ws.onerror = (e) => {
        console.log('websocket error')
        console.log(e)
    }
    ws.onclose = () => {
        console.log(`quiting room: ${roomname}`)
    }
}

createWebSocket(wsurl)
setWebSocketCallback(ws)

const parseMessage = (msg) => {
    let ret = null
    try {
        let parsed = JSON.parse(msg)
        let color = parsed.color
        if (!('color' in parsed)) {
            color = 'Gray'
        }
        ret = {
            username: parsed.username,
            message: parsed.message,
            color: color
        }
    } catch (e) {
        ret = {
            username: '- -',
            message: msg,
            color: 'Gray'
        }
    }

    return ret
}

const constructMsg = (username, msg) => {
    let obj = {
        username: username,
        message: msg,
        color: userColorTheme
    }
    return JSON.stringify(obj)
}

const sendMessage = (msg, ws) => {
    let wholeMsg = constructMsg(username, msg)
    ws.send(wholeMsg)
    inputText.value = ''
}

const quitChatRoom = () => {
    ws.close()
    ws = null
    router.push('/chatrooms')
}

const refreshMsg = () => {
    getMessages()
}

var now = function () {
    var date = new Date(new Date().getTime() + (parseInt(new Date().getTimezoneOffset() / 60) + 8) * 3600 * 1000).toString();
    var dateArr = date.split(" ");
    return dateArr[1] + " " + dateArr[2] + " " + dateArr[3] + " " + dateArr[4];
}

</script>
<template>
    <div class="bg-dark chat-wrap">
        <div>
            <h3 class="text-white text-center">{{ roomname }}</h3>
            <div>
                <button class="me-1 btn btn-danger btn-sm" @click="quitChatRoom()">Quit</button>
                <button class="mx-1 btn btn-info btn-sm" @click="refreshMsg()">Refresh</button>
                <button class="mx-1 btn btn-primary btn-sm" style="width:80px"
                    @click="sendMessage(inputText, ws)">Send</button>
            </div>
        </div>
        <div class="mt-2">
            <textarea class="w-100" type="text" rows="2" v-model="inputText"
                @keyup.enter="sendMessage(inputText, ws)"></textarea>
        </div>
        <div class="mt-2 table-responsive">
            <table class="table table-hover table-borderless">
                <tbody>
                    <tr v-for="message in messageData">
                        <td>
                            <div class="text-center text-white float-start">
                                <div class="avatar border border-white"
                                    :style="{ 'background-color': convertSelectToColor(message.color) }">
                                    {{ message.username.charAt(0) }}
                                </div>
                                <div>
                                    {{ message.username }}
                                </div>
                            </div>
                            <div class="bub float-start px-3 ms-2 border border-white"
                                :style="{ 'background-color': convertSelectToColor(message.color) }">
                                {{ message.message }}
                            </div>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
    <!-- <div class="col-sm-2"></div> -->
</template>
<style>
/* .chat-wrap {
    min-height: 100vh;
} */

.avatar {
    /* background-color: gray; */
    line-height: 50px;
    width: 50px;
    font-size: 30px;
    /* display: block; */
    border-radius: 5px;
}

.username {
    width: 50px;
}

.bub {
    word-break: break-all;
    font-family: Microsoft Yahei;
    position: relative;
    /* width: 200px; */
    /* height: 50px; */
    color: #fff;
    /* font-size: 1em; */
    font-size: 1rem;
    line-height: 50px;
    text-align: start;
    /* border: 1px solid teal; */
    border-radius: 5px;
    /* background: teal; */
    /* background: gray; */
    max-width: 70vw;
}

.bub::after {
    content: '';
    position: absolute;
    width: 0;
    height: 0;
    /* 箭头靠右边 */
    /* top: 13px;
    right: -10px;
    border-top: 10px solid transparent;
    border-bottom: 10px solid transparent;
    border-left: 10px solid teal; */
    /* 箭头靠下边 */
    /* left: 25px;
            bottom: -10px;
            border-left: 10px solid transparent;
            border-right: 10px solid transparent;
            border-top: 10px solid teal; */
    /* 箭头靠左边 */
    top: 13px;
    left: -6px;
    border-top: 6px solid transparent;
    border-bottom: 6px solid transparent;
    border-right: 6px solid white;
    /* 箭头靠下边 */
    /* top: -10px;
            left: 25px;
            border-left: 10px solid transparent;
            border-right: 10px solid transparent;
            border-bottom: 10px solid teal; */
}
</style>
