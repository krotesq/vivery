import { z } from "zod"
import type { ComposerTranslation } from "vue-i18n"

export function createStreamTargetSchema(t: ComposerTranslation) {
  return z.object({
    name: z.string().min(1, t("validation.mandatory")),
    streamLink: z.string().min(1, t("validation.mandatory")),
    streamKey: z.string().min(1, t("validation.mandatory")),
  })
}
