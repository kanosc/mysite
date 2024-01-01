<script setup>
import { useRouter } from 'vue-router'
import { ref } from 'vue'
import { convertSelectToColor } from '../common/common'
const username = ref('')
const selectedColor = ref('Gray')

const router = useRouter()
// const props = defineProps(['option', 'info'])
const loginUser = () => {
    let purename = username.value.trim()
    if (purename !== "") {
        localStorage.setItem("username-chat", purename)
        localStorage.setItem("color-theme", selectedColor.value)
        router.push('/chatrooms')
    } else {
        alert("invalid username")
    }
}
</script>
<template>
    <div class="login-wrap">
        <main class="form-signin w-100 mx-auto">
            <form action="/chatrooms" method="post" enctype="multipart/form-data" @submit.prevent="loginUser()">
                <h1 class="h3 mb-3 fw-normal">Login</h1>
                <!-- <label for="inputid" class="form-label">{{ props.option }}</label> -->
                <label for="inputid" class="form-label">username</label>
                <input type="text" id="inputid" class="form-control" aria-describedby="usernameHelpBlock"
                    v-model="username">
                <!-- <div id="usernameHelpBlock" class="form-text">
                     {{ props.info }}
                    please input your username
                </div> -->
                <label for="themeid" class="form-label">
                    colortheme
                </label>
                <select class="form-select" id="themeid" v-model="selectedColor">
                    <option value="Dark">Dark</option>
                    <option value="Red">Red</option>
                    <option value="Yellow">Yellow</option>
                    <option value="Green">Green</option>
                    <option value="Orange">Orange</option>
                    <option value="Blue">Blue</option>
                    <option value="Pink">Pink</option>
                    <option value="Gray">Gray(Default)</option>
                </select>
                <div class="mt-3">
                    <span class="colorPrev" style=" width: 4rem;height: 30px;"
                        :style="{ 'background-color': convertSelectToColor(selectedColor), 'color': convertSelectToColor(selectedColor) }">
                        ab
                    </span>
                    {{ selectedColor }}
                </div>

                <button class="w-100 btn btn-lg btn-dark mt-3" type="submit">Login</button>
                <p class="mt-3 mb-3 text-muted">&copy; 2023–2028</p>
            </form>
        </main>
    </div>
</template>
<style>
.colorPrev {
    background-color: #333333;
    color: #333333;
}

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