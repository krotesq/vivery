import { ref } from "vue"
import { useI18n } from "vue-i18n"
import { toast } from "vue-sonner"
import type { StreamTarget, StreamTargetInput, StreamKind } from "@/types"

function makeStreamTarget(input: StreamTargetInput): StreamTarget {
  return {
    id: crypto.randomUUID(),
    name: input.name,
    streamLink: input.streamLink,
    streamKey: input.streamKey,
  }
}

// Module-scoped so state survives navigation within the SPA (mock only).
const sources = ref<StreamTarget[]>([])
const targets = ref<StreamTarget[]>([])

export function useStreamTargets() {
  const { t } = useI18n()

  function listFor(kind: StreamKind) {
    return kind === "source" ? sources : targets
  }

  function create(kind: StreamKind, input: StreamTargetInput) {
    // TODO: replace mock with streamingService.create{Source|Target}(input)
    listFor(kind).value.push(makeStreamTarget(input))
    toast.success(t("streaming.created"))
  }

  function update(kind: StreamKind, id: string, input: StreamTargetInput) {
    // TODO: replace mock with streamingService.update{Source|Target}(id, input)
    const list = listFor(kind)
    const idx = list.value.findIndex((e) => e.id === id)
    if (idx !== -1) {
      list.value[idx] = { ...list.value[idx], ...input }
      toast.success(t("streaming.updated"))
    }
  }

  function remove(kind: StreamKind, id: string) {
    // TODO: replace mock with streamingService.delete{Source|Target}(id)
    const list = listFor(kind)
    list.value = list.value.filter((e) => e.id !== id)
    toast.success(t("streaming.deleted"))
  }

  return { sources, targets, create, update, remove }
}
