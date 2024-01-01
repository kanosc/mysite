<script setup>
import { ref } from 'vue'
const imgAddr = 'https://t.mwm.moe/fj/'
const show = ref(false)
const url = ref(imgAddr)

function loaded() {
    console.log("image loaded")
    show.value = true
}

function refresh() {
    location.reload()
}

</script>
<template>
    <!-- <h1 class="text-center">
        fortune card
    </h1>
    <img class="img-fluid" decoding="async" src="https://t.mwm.moe/fj/"> -->
    <div class="container">
        <div v-show="show" class="card mx-auto mt-3" style="max-width: 80rem;">
            <img @load="loaded" :src="url" class="card-img-top" alt="press [New image] button to get image"
                decoding="async">
            <div class="card-body">
                <p class="card-text">Later, respectively, wander and suffer sorrow.</p>
                <button @click="refresh" class="btn btn-primary float-end">Random</button>
            </div>
        </div>
        <div v-if="!show">
            <h5 v-if="!show">Resource Loading...</h5>
            <div v-if="!show" class="spinner-border text-info"></div>
        </div>
    </div>
</template>

