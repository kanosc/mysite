<script setup>
import { ref } from 'vue';

const imgAddr = ref("")
const defaultAddr = ref("https://w.wallhaven.cc/full/gp/wallhaven-gp1yx7.jpg")
let clientImage = localStorage.getItem('image-addr')
if (clientImage != null) {
    defaultAddr.value = clientImage
}

const setImageAddr = () => {
    defaultAddr.value = imgAddr.value
    localStorage.setItem('image-addr', imgAddr.value)
}

</script>
<template>
    <div class="container">
        <div class="mt-3 mb-3 row">
            <div class="col">
                <input type="text" class="form-control float-start me-2" v-model="imgAddr"
                    placeholder="input the address of image. exp: http://xxxx/x.jpg">
            </div>
            <div class="col">
                <button class="btn btn-primary" @click="setImageAddr()">Set Image</button>
            </div>
        </div>
        <div>
            <img decoding="async" class="img-fluid" :src="defaultAddr" alt="Image address not set">
        </div>
    </div>
</template>
