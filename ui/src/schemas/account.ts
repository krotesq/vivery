import { z } from 'zod'
import type { ComposerTranslation } from 'vue-i18n'

export function createLoginSchema(t: ComposerTranslation) {
  return z.object({
    username: z.string().min(1, t('validation.mandatory')),
    password: z.string().min(1, t('validation.mandatory')),
  })
}

export function createProfileSchema(t: ComposerTranslation) {
  return z
    .object({
      username: z.string().min(1, t('validation.mandatory')),
      password: z.string().optional(),
      confirmPassword: z.string().optional(),
    })
    .refine((data) => (data.password ?? '') === (data.confirmPassword ?? ''), {
      message: t('validation.passwordsMismatch'),
      path: ['confirmPassword'],
    })
}