<script setup>
import { ref } from 'vue'
const roomDetail = ref(null)
const emit = defineEmits(['create-room-success'])

const createRoom = async () => {
    let response = await fetch('/api/chat/v1/room', {
        method: 'POST',
        body: new FormData(roomDetail.value)
    })
    let rspText = await response.text()
    if (response.status == 200) {
        alert('room create success')
        emit('create-room-success')
    } else {
        alert(rspText)
    }
}

</script>
<template>
    <button class="btn btn-info mt-2 float-start" data-bs-toggle="modal" data-bs-target="#create-room">
        create room
    </button>
    <div class="modal" id="create-room">
        <div class="modal-dialog">
            <div class="modal-content">
                <!-- 模态框头部 -->
                <div class="modal-header">
                    <span class="modal-title">creating new room</span>
                    <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                </div>
                <!-- 模态框内容 -->
                <div class="modal-body">
                    <form action="/create_room" method="POST" enctype="multipart/form-data" class="was-validated"
                        @submit.prevent="createRoom()" ref="roomDetail">
                        <input class="form-control" type="text" name="roomname" placeholder="input room name" required />
                        <div class="valid-feedback">ok</div>
                        <div class="invalid-feedback">please input room name</div>
                        <hr>
                        <input class="form-control" type="text" name="password" placeholder="input password(optinal)" />
                        <hr>
                        <input class="form-control" type="text" name="maxuser" placeholder="input max users(optrinal)" />
                        <!-- <input type="file" class="select-file" multiple="multiple" id="upload" name="upload"/> -->
                        <hr>
                        <button type="submit" class="btn btn-primary">Submit</button>
                    </form>
                </div>
                <!-- 模态框底部 -->
                <div class="modal-footer">
                    <button type="button" class="btn btn-danger" data-bs-dismiss="modal">close</button>
                </div>

            </div>
        </div>
    </div>
</template>