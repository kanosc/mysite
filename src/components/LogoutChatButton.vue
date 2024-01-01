<script setup>
import { useRouter } from 'vue-router'
const router = useRouter()
const deleteUsername = () => {
    let username = localStorage.getItem('username-chat')
    if (username != null) {
        localStorage.removeItem('username-chat')
    }
    router.push('/login/username')
}
</script>
<template>
    <button class="btn btn-danger mt-2 float-start mx-2" @click="deleteUsername()">
        log out
    </button>
</template>