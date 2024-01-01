<script setup>
import { ref } from 'vue'
import { copyList } from '../common/common.js'
import CreateRoomButton from '../components/CreateRoomButton.vue'
import LogoutChatButton from '../components/LogoutChatButton.vue'

const username = localStorage.getItem('username-chat')
let helloMsg = 'hello'
if (username != null) {
    helloMsg = 'hello ' + username
}

let originRooms = null
const roomList = ref([])
const getRooms = async () => {
    try {
        const res = await fetch("/api/chat/v1/rooms")
        const rsp = await res.json()
        originRooms = rsp.rooms
        roomList.value = copyList(originRooms)
    } catch (e) {
        alert(e)
    }
}
getRooms()

const makeRoomeAddr = (roomname) => {
    return "/page/room/" + roomname
}
</script>
<template>
    <div class="container-fluid">
        <div class="row">
            <div class="col-sm-2"></div>
            <div class="col-sm-8">
                <div class="mt-3 container-fluid">
                    <CreateRoomButton @create-room-success="getRooms()" />
                    <LogoutChatButton />
                </div>
                <div class="table-responsive container-fluid">
                    <table class="table table-hover">
                        <thead>
                            <tr>
                                <th>Room Name</th>
                                <th><span>Operation</span></th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="room in roomList">
                                <td>
                                    <span>{{ room }}</span>
                                </td>
                                <td>
                                    <span>
                                        <AppLink :to="makeRoomeAddr(room)"><button class="btn btn-warning">进入房间</button>
                                        </AppLink>
                                    </span>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
            <div class="col-sm-2"></div>
        </div>
    </div>
</template>