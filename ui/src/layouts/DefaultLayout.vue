<script setup lang="ts">
import { computed, type Component } from "vue"
import { useRoute, useRouter } from "vue-router"
import {
  HomeIcon,
  ArrowsRightLeftIcon,
  Squares2X2Icon,
  UserCircleIcon
} from "@heroicons/vue/24/outline"
import LiquidGlassNav, { type GlassNavItem } from "@/components/LiquidGlassNav.vue"

const route = useRoute()
const router = useRouter()

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
</script>

<template>
  <div class="top-bar-right">
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