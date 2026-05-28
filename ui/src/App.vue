<script setup lang="ts">
import { Toaster } from 'vue-sonner'
import { useThemeStore } from '@/stores/theme'
useThemeStore()
</script>

<template>
  <Toaster />
  <RouterView/>
</template>
