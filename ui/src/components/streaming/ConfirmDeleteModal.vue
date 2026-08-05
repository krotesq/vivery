<script setup lang="ts">
import { ref, watch } from "vue"

const props = defineProps<{ open: boolean; title: string }>()
const emit = defineEmits<{ confirm: []; cancel: [] }>()

const dialog = ref<HTMLDialogElement | null>(null)

watch(
  () => props.open,
  (open) => {
    if (open) dialog.value?.showModal()
    else dialog.value?.close()
  },
)
</script>

<template>
  <dialog ref="dialog" class="modal" @close="emit('cancel')">
    <div class="modal-box">
      <h3 class="text-lg font-semibold">{{ $t("streaming.confirmDelete", { title }) }}</h3>
      <div class="modal-action">
        <button type="button" class="btn btn-ghost" @click="emit('cancel')">
          {{ $t("streaming.cancel") }}
        </button>
        <button type="button" class="btn btn-error" @click="emit('confirm')">
          {{ $t("streaming.delete") }}
        </button>
      </div>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button>close</button>
    </form>
  </dialog>
</template>
