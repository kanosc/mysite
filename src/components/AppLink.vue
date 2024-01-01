<template>
    <a v-if="isExternal" target="_blank" rel="noopener" :href="to">
        <slot />
    </a>
    <router-link v-else v-bind="$props">
        <slot />
    </router-link>
</template>
<script setup>
import { RouterLink } from 'vue-router'
import { computed } from 'vue'
const props = defineProps({ ...RouterLink.props })
const isExternal = computed(() => typeof props.to === 'string' && props.to.startsWith('http'))
</script>