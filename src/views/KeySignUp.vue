<script setup>
import { useRouter } from 'vue-router'
import { ref } from 'vue'
const password = ref('')
const password2 = ref('')
const signupForm = ref(null)

const router = useRouter()
const keySignUp = async () => {
    alert("sign up success")
    // console.log(password)
    // try {
    //     let response = await fetch('/api/chat/v1/auth/' + props.roomname, {
    //         method: 'POST',
    //         body: new FormData(pwdForm.value)
    //     })
    //     if (response.status != 200) {
    //         alert(await response.text())
    //         return
    //     }
    //     let authRet = await response.json()
    //     let token = authRet.token
    //     console.log(`received token ${token}`)
    //     localStorage.setItem(props.roomname + '-token', token)
    //     router.push(`/page/room/${props.roomname}`)

    // } catch (e) {
    //     console.log(e)
    // }

}
</script>
<template>
    <div class="login-wrap">
        <main class="form-signin w-100 mx-auto">
            <form action="/#" method="post" @submit.prevent="keySignUp()" ref="signupForm">
                <h1 class="h3 mb-3 fw-normal">Sign Up(注册)</h1>
                <!-- <label for="inputid" class="form-label">{{ props.option }}</label> -->
                <label for="t-name" class="form-label">email</label>
                <input type="email" id="t-name" class="form-control" name="email" v-model="username"
                    placeholder="name@example.com">
                <label for="t-pwd" class="form-label">password for new account</label>
                <input type="password" id="t-pwd" class="form-control" name="password" v-model="password">
                <label for="t-pwd2" class="form-label">please input password again</label>
                <input type="password" id="t-pwd2" class="form-control" name="password2" v-model="password2">
                <span class="text-danger">Notice: system will not collect your password, please remember or put it in secure
                    place.</span>

                <button class="w-100 btn btn-lg btn-dark mt-3" type="submit">Sign Up</button>
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