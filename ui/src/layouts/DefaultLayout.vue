<script setup lang="ts">
import { computed, type Component } from "vue"
import { useRoute, useRouter } from "vue-router"
import {
  HomeIcon,
  ArrowsRightLeftIcon,
  Squares2X2Icon,
  UserCircleIcon,
  ArrowRightStartOnRectangleIcon,
  SunIcon,
  MoonIcon,
} from "@heroicons/vue/24/outline"
import LiquidGlassNav, { type GlassNavItem } from "@/components/LiquidGlassNav.vue"
import { useAuthStore } from "@/stores/auth"
import { useThemeStore } from "@/stores/theme"

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const themeStore = useThemeStore()

type NavItem = GlassNavItem & { icon: Component }

const items: NavItem[] = [
  { value: "dashboard", label: "Home", icon: HomeIcon },
  { value: "flow-editor", label: "Flow", icon: ArrowsRightLeftIcon },
  { value: "source-target", label: "Sources", icon: Squares2X2Icon },
  { value: "profil", label: "Profil", icon: UserCircleIcon },
]

const current = computed<string | number>({
  get: () => (route.name as string) ?? "dashboard",
  set: (name) => {
    router.push({ name: String(name) })
  },
})

async function onLogout() {
  await auth.logout()
  await router.push({ name: "login" })
}

function toggleTheme() {
  themeStore.setTheme(themeStore.theme === "dark" ? "light" : "dark")
}
</script>

<template>
  <button class="logout-btn btn btn-outline btn-error" @click="onLogout">
    <ArrowRightStartOnRectangleIcon class="size-5" />
    {{ $t("account.logout") }}
  </button>

  <div class="top-bar-right">
    <label class="swap swap-rotate theme-toggle">
      <input
        type="checkbox"
        :checked="themeStore.theme === 'dark'"
        aria-label="Toggle theme"
        @change="toggleTheme"
      />
      <SunIcon class="swap-on size-6" />
      <MoonIcon class="swap-off size-6" />
    </label>

    <div class="brand-badge bg-base-200">Vivery</div>
  </div>

  <nav class="default-nav">
    <LiquidGlassNav v-model="current" :items="items" aria-label="Main navigation">
      <template #icon="{ item }">
        <component :is="(item as NavItem).icon" :title="item.label" />
      </template>
    </LiquidGlassNav>
  </nav>
  <main class="page-content min-h-screen bg-base-300 text-base-content">
    <RouterView />
  </main>
</template>

<style scoped>
.default-nav {
  position: fixed;
  z-index: 50;
  bottom: 40px;
  left: 50%;
  translate: -50%;
}

.logout-btn {
  position: fixed;
  z-index: 50;
  top: 12px;
  left: 20px;
}

.top-bar-right {
  position: fixed;
  z-index: 50;
  top: 0;
  right: 0;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.theme-toggle {
  cursor: pointer;
}

.brand-badge {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 56px;
  padding: 0 1.5rem;
  font-weight: 600;
  border-bottom-left-radius: 30px;
}

.page-content {
  padding: 76px 20px 20px;
}
</style>
