<script setup lang="ts">
import { ref } from "vue"
import { useI18n } from "vue-i18n"
import StreamTargetSection from "@/components/streaming/StreamTargetSection.vue"
import StreamTargetFormModal from "@/components/streaming/StreamTargetFormModal.vue"
import ConfirmDeleteModal from "@/components/streaming/ConfirmDeleteModal.vue"
import { useStreamTargets } from "@/composables/useStreamTargets"
import type { StreamTarget, StreamTargetInput, StreamKind } from "@/types"

const { t } = useI18n()
const { sources, targets, create, update, remove } = useStreamTargets()

const formOpen = ref(false)
const confirmOpen = ref(false)
const activeKind = ref<StreamKind>("source")
const formMode = ref<"add" | "edit">("add")
const editing = ref<StreamTarget | null>(null)
const deleting = ref<StreamTarget | null>(null)

function openAdd(kind: StreamKind) {
  activeKind.value = kind
  formMode.value = "add"
  editing.value = null
  formOpen.value = true
}

function openEdit(kind: StreamKind, item: StreamTarget) {
  activeKind.value = kind
  formMode.value = "edit"
  editing.value = item
  formOpen.value = true
}

function openDelete(kind: StreamKind, item: StreamTarget) {
  activeKind.value = kind
  deleting.value = item
  confirmOpen.value = true
}

function onSubmit(values: StreamTargetInput) {
  if (formMode.value === "edit" && editing.value) {
    update(activeKind.value, editing.value.id, values)
  } else {
    create(activeKind.value, values)
  }
  formOpen.value = false
}

function onFormDelete() {
  deleting.value = editing.value
  confirmOpen.value = true
}

function onConfirmDelete() {
  if (deleting.value) remove(activeKind.value, deleting.value.id)
  confirmOpen.value = false
  deleting.value = null
  formOpen.value = false
}
</script>

<template>
  <div>
    <section class="flex min-h-[calc(100vh-160px)] flex-col gap-12 rounded-box bg-base-100 p-6 shadow-lg sm:p-10">
      <StreamTargetSection
        :title="t('streaming.source')"
        :items="sources"
        @add="openAdd('source')"
        @edit="(item) => openEdit('source', item)"
        @delete="(item) => openDelete('source', item)"
      />

      <StreamTargetSection
        :title="t('streaming.target')"
        :items="targets"
        @add="openAdd('target')"
        @edit="(item) => openEdit('target', item)"
        @delete="(item) => openDelete('target', item)"
      />
    </section>

    <StreamTargetFormModal
      :open="formOpen"
      :kind="activeKind"
      :mode="formMode"
      :initial="editing"
      @submit="onSubmit"
      @delete="onFormDelete"
      @close="formOpen = false"
    />

    <ConfirmDeleteModal
      :open="confirmOpen"
      :title="deleting?.name ?? ''"
      @confirm="onConfirmDelete"
      @cancel="confirmOpen = false"
    />
  </div>
</template>
