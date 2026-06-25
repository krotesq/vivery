import { createApp } from 'vue'
import './style.css'
import 'vue-sonner/style.css'
import App from './App.vue'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import { VueQueryPlugin } from '@tanstack/vue-query'
import router from './router'
import i18n from "./i18n"
import { initializeTheme, useThemeStore } from '@/stores/theme'

const app = createApp(App)

initializeTheme()

const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

app.use(pinia)
useThemeStore()
app.use(router)
app.use(VueQueryPlugin)
app.use(i18n)

app.mount("#app")
