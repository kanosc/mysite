import { createApp } from 'vue'
import router from './router'
import AppLink from './components/AppLink.vue'
import App from './App.vue'

createApp(App)
    .component('AppLink', AppLink)
    .use(router)
    .mount('#app')
