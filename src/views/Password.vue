<script setup>
import { useRouter } from 'vue-router'
import { ref } from 'vue'
const password = ref('')
const pwdForm = ref(null)

const router = useRouter()
const props = defineProps(['roomname'])
const loginPassword = async () => {
    console.log(password)
    try {
        let response = await fetch('/api/chat/v1/auth/' + props.roomname, {
            method: 'POST',
            body: new FormData(pwdForm.value)
        })
        if (response.status != 200) {
            alert(await response.text())
            return
        }
        let authRet = await response.json()
        let token = authRet.token
        console.log(`received token ${token}`)
        localStorage.setItem(props.roomname + '-token', token)
        router.push(`/page/room/${props.roomname}`)

    } catch (e) {
        console.log(e)
    }

}
</script>
<template>
    <div class="login-wrap">
        <main class="form-signin w-100 mx-auto">
            <form action="/#" method="post" @submit.prevent="loginPassword()" ref="pwdForm">
                <h1 class="h3 mb-3 fw-normal">Login</h1>
                <!-- <label for="inputid" class="form-label">{{ props.option }}</label> -->
                <label for="inputid" class="form-label">password for room {{ props.roomname }}</label>
                <input type="password" id="inputid" class="form-control" name="password" v-model="password">

                <button class="w-100 btn btn-lg btn-dark mt-3" type="submit">Login</button>
                <p class="mt-3 mb-3 text-muted">&copy; 2023–2028</p>
            </form>
        </main>
    </div>
</template>
<style>
.login-wrap {
    display: flex;
    height: 60vh;
    padding-top: 40px;
    padding-bottom: 40px;
    align-items: center;
}

.form-signin {
    max-width: 330px;
    padding: 15px;
}
</style>