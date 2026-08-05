<script setup lang="ts">
import { ref, computed } from "vue"
import { PlusIcon, MagnifyingGlassIcon } from "@heroicons/vue/24/outline"
import StreamTargetCard from "./StreamTargetCard.vue"
import type { StreamTarget } from "@/types"

const props = defineProps<{ title: string; items: StreamTarget[] }>()
const emit = defineEmits<{
  add: []
  edit: [item: StreamTarget]
  delete: [item: StreamTarget]
}>()

const search = ref("")
const filtered = computed(() => {
  const q = search.value.trim().toLowerCase()
  if (!q) return props.items
  return props.items.filter((i) => i.name.toLowerCase().includes(q))
})
</script>

<template>
  <section class="flex flex-col gap-4">
    <h2 class="text-2xl font-semibold">{{ title }}</h2>

    <div class="flex items-stretch gap-3">
      <button type="button" class="btn btn-square h-14 w-14 text-primary" :aria-label="title" @click="emit('add')">
        <PlusIcon class="size-6" />
      </button>
      <label class="input h-14 flex-1">
        <MagnifyingGlassIcon class="size-4 opacity-50" />
        <input v-model="search" type="text" :placeholder="$t('streaming.searchPlaceholder')" />
      </label>
    </div>

    <div class="flex min-h-40 flex-wrap gap-4 rounded-box bg-base-300 p-4">
      <StreamTargetCard
        v-for="item in filtered"
        :key="item.id"
        :item="item"
        @edit="emit('edit', item)"
        @delete="emit('delete', item)"
      />
    </div>
  </section>
</template>
