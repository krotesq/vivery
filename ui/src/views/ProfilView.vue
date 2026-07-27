<script setup lang="ts">
import { ref } from "vue"
import { useI18n } from "vue-i18n"
import { useForm } from "vee-validate"
import { toTypedSchema } from "@vee-validate/zod"
import { toast } from "vue-sonner"
import { BookmarkIcon } from "@heroicons/vue/24/outline"
import { createProfileSchema } from "@/schemas/account"
import { useAuthStore } from "@/stores/auth"

const auth = useAuthStore()
const { t } = useI18n()

const { handleSubmit, defineField, meta, errors } = useForm({
  validationSchema: toTypedSchema(createProfileSchema(t)),
  initialValues: {
    username: auth.user?.username ?? "",
    password: "",
    confirmPassword: "",
  },
})

const [username, usernameAttrs] = defineField("username")
const [password, passwordAttrs] = defineField("password")
const [confirmPassword, confirmPasswordAttrs] = defineField("confirmPassword")

const loading = ref(false)

const onSubmit = handleSubmit(
  async () => {
    loading.value = true
    try {
      toast.success(t("account.saved"))
      password.value = ""
      confirmPassword.value = ""
    } finally {
      loading.value = false
    }
  },
  ({ errors }) => {
    const firstError = Object.values(errors)[0]
    if (firstError) toast.error(firstError)
  },
)
</script>

<template>
  <section class="profile-card flex flex-col rounded-box bg-base-100 p-6 shadow-lg sm:p-10">
    <form @submit.prevent="onSubmit" class="flex flex-1 flex-col" aria-label="Profil">
      <div class="flex w-full max-w-sm flex-col gap-6">
        <div class="flex flex-col gap-2">
          <label for="profile-username" class="font-medium">{{ $t("account.username") }}</label>
          <input
            id="profile-username"
            v-model="username"
            v-bind="usernameAttrs"
            type="text"
            :placeholder="$t('account.usernamePlaceholder')"
            autocomplete="username"
            class="input w-full"
            :class="{ 'input-error': errors.username }"
          />
          <span v-if="errors.username" class="text-sm text-error">{{ errors.username }}</span>
        </div>

        <div class="flex flex-col gap-2">
          <label for="profile-password" class="font-medium">{{ $t("account.password") }}</label>
          <input
            id="profile-password"
            v-model="password"
            v-bind="passwordAttrs"
            type="password"
            :placeholder="$t('account.passwordPlaceholder')"
            autocomplete="new-password"
            class="input w-full"
            :class="{ 'input-error': errors.password }"
          />
          <span v-if="errors.password" class="text-sm text-error">{{ errors.password }}</span>
        </div>

        <div class="flex flex-col gap-2">
          <label for="profile-password-confirm" class="font-medium">{{ $t("account.confirmPassword") }}</label>
          <input
            id="profile-password-confirm"
            v-model="confirmPassword"
            v-bind="confirmPasswordAttrs"
            type="password"
            :placeholder="$t('account.confirmPasswordPlaceholder')"
            autocomplete="new-password"
            class="input w-full"
            :class="{ 'input-error': errors.confirmPassword }"
          />
          <span v-if="errors.confirmPassword" class="text-sm text-error">{{ errors.confirmPassword }}</span>
        </div>
      </div>

      <div class="mt-auto flex justify-end">
        <button type="submit" :disabled="loading || !meta.dirty || !meta.valid" class="btn btn-success">
          <BookmarkIcon class="size-5" />
          {{ $t("account.save") }}
        </button>
      </div>
    </form>
  </section>
</template>

<style scoped>
.profile-card {
  height: calc(100vh - 260px);
}
</style>
