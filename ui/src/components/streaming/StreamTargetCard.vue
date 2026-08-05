<script setup lang="ts">
import { SignalIcon, EllipsisHorizontalIcon } from "@heroicons/vue/24/outline"
import type { StreamTarget } from "@/types"

defineProps<{ item: StreamTarget }>()
const emit = defineEmits<{ edit: []; delete: [] }>()
</script>

<template>
  <div class="card card-border relative w-72 border-base-content/10 bg-base-200/40">
    <div class="dropdown dropdown-end absolute right-2 top-2">
      <div
        tabindex="0"
        role="button"
        class="btn btn-ghost btn-xs btn-square"
        :aria-label="$t('streaming.edit')"
      >
        <EllipsisHorizontalIcon class="size-5" />
      </div>
      <ul tabindex="0" class="dropdown-content menu z-10 w-36 rounded-box bg-base-100 p-2 shadow">
        <li><a @click="emit('edit')">{{ $t("streaming.edit") }}</a></li>
        <li><a class="text-error" @click="emit('delete')">{{ $t("streaming.delete") }}</a></li>
      </ul>
    </div>

    <button type="button" class="card-body flex-row items-center gap-4 text-left" @click="emit('edit')">
      <SignalIcon class="size-10 shrink-0 opacity-40" />
      <span class="truncate font-medium">{{ item.name }}</span>
    </button>
  </div>
</template>
