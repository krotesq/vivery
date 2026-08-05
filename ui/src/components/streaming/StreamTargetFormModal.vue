<script setup lang="ts">
import { ref, watch, computed } from "vue"
import { useI18n } from "vue-i18n"
import { useForm } from "vee-validate"
import { toTypedSchema } from "@vee-validate/zod"
import { toast } from "vue-sonner"
import {
  Square3Stack3DIcon,
  XMarkIcon,
  LinkIcon,
  LockClosedIcon,
  TrashIcon,
} from "@heroicons/vue/24/outline"
import { createStreamTargetSchema } from "@/schemas/streaming"
import type { StreamTarget, StreamTargetInput, StreamKind } from "@/types"

const props = defineProps<{
  open: boolean
  kind: StreamKind
  mode: "add" | "edit"
  initial?: StreamTarget | null
}>()
const emit = defineEmits<{ submit: [values: StreamTargetInput]; close: []; delete: [] }>()

const { t } = useI18n()
const dialog = ref<HTMLDialogElement | null>(null)

const { handleSubmit, defineField, errors, resetForm } = useForm({
  validationSchema: toTypedSchema(createStreamTargetSchema(t)),
  initialValues: { name: "", streamLink: "", streamKey: "" },
})

const [name, nameAttrs] = defineField("name")
const [streamLink, streamLinkAttrs] = defineField("streamLink")
const [streamKey, streamKeyAttrs] = defineField("streamKey")

const titleKey = computed(() => {
  if (props.kind === "source") return props.mode === "add" ? "streaming.addSource" : "streaming.editSource"
  return props.mode === "add" ? "streaming.addTarget" : "streaming.editTarget"
})

watch(
  () => props.open,
  (open) => {
    if (open) {
      resetForm({
        values: {
          name: props.initial?.name ?? "",
          streamLink: props.initial?.streamLink ?? "",
          streamKey: props.initial?.streamKey ?? "",
        },
      })
      dialog.value?.showModal()
    } else {
      dialog.value?.close()
    }
  },
)

const onSubmit = handleSubmit(
  (values) => emit("submit", values),
  ({ errors }) => {
    const first = Object.values(errors)[0]
    if (first) toast.error(first)
  },
)
</script>

<template>
  <dialog ref="dialog" class="modal" @close="emit('close')">
    <div class="modal-box max-w-2xl">
      <div class="mb-6 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <Square3Stack3DIcon class="size-6 text-primary" />
          <h3 class="text-lg font-semibold">{{ $t(titleKey) }}</h3>
        </div>
        <button type="button" class="btn btn-ghost btn-sm btn-square" :aria-label="$t('streaming.cancel')" @click="emit('close')">
          <XMarkIcon class="size-5" />
        </button>
      </div>

      <form class="flex flex-col gap-6" @submit.prevent="onSubmit">
        <div class="flex flex-col gap-2">
          <label for="stream-target-name" class="font-medium">{{ $t("streaming.name") }}</label>
          <label class="input w-full">
            <input
              id="stream-target-name"
              v-model="name"
              v-bind="nameAttrs"
              type="text"
              :placeholder="$t('streaming.namePlaceholder')"
            />
          </label>
          <span v-if="errors.name" class="text-sm text-error">{{ errors.name }}</span>
        </div>

        <div class="flex flex-col gap-6 sm:flex-row">
          <div class="flex flex-1 flex-col gap-2">
            <label for="stream-target-link" class="font-medium">{{ $t("streaming.streamLink") }}</label>
            <label class="input w-full">
              <LinkIcon class="size-4 opacity-50" />
              <input
                id="stream-target-link"
                v-model="streamLink"
                v-bind="streamLinkAttrs"
                type="text"
                :placeholder="$t('streaming.streamLinkPlaceholder')"
              />
            </label>
            <span v-if="errors.streamLink" class="text-sm text-error">{{ errors.streamLink }}</span>
          </div>

          <div class="flex flex-1 flex-col gap-2">
            <label for="stream-target-key" class="font-medium">{{ $t("streaming.streamKey") }}</label>
            <label class="input w-full">
              <LockClosedIcon class="size-4 opacity-50" />
              <input
                id="stream-target-key"
                v-model="streamKey"
                v-bind="streamKeyAttrs"
                type="text"
                :placeholder="$t('streaming.streamKeyPlaceholder')"
              />
            </label>
            <span v-if="errors.streamKey" class="text-sm text-error">{{ errors.streamKey }}</span>
          </div>
        </div>

        <div class="mt-6 flex items-center gap-2">
          <button
            v-if="mode === 'edit'"
            type="button"
            class="btn btn-outline btn-error"
            @click="emit('delete')"
          >
            <TrashIcon class="size-5" />
            {{ $t("streaming.delete") }}
          </button>
          <div class="ml-auto flex gap-2">
            <button type="button" class="btn btn-ghost" @click="emit('close')">
              {{ $t("streaming.cancel") }}
            </button>
            <button type="submit" class="btn btn-primary">
              {{ $t("streaming.save") }}
            </button>
          </div>
        </div>
      </form>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button>close</button>
    </form>
  </dialog>
</template>
