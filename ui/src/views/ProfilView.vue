<script setup lang="ts">
import { useRouter } from "vue-router"
import { useThemeStore } from "@/stores/theme"
import { ref } from "vue"
import { useI18n } from "vue-i18n"
import { useForm } from "vee-validate"
import { toTypedSchema } from "@vee-validate/zod"
import { toast } from "vue-sonner"
import { 
  ArrowRightStartOnRectangleIcon,
  SunIcon,
  MoonIcon 
} from "@heroicons/vue/24/outline"
import { createProfileSchema } from "@/schemas/account"
import { useAuthStore } from "@/stores/auth"

const auth = useAuthStore()
const router = useRouter()
const themeStore = useThemeStore()
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

async function onLogout() {
  await auth.logout()
  await router.push({ name: "login" })
}

function toggleTheme() {
  themeStore.setTheme(themeStore.theme === "dark" ? "emerald" : "dark")
}
</script>

<template>
  <section class="profile-card flex flex-col rounded-box bg-base-100 p-6 shadow-lg sm:p-10">
    <div class="flex justify-between mb-10">
      <label class="swap swap-rotate theme-toggle">
        <input
          type="checkbox"
          :checked="themeStore.theme === 'dark'"
          aria-label="Toggle theme"
          @change="toggleTheme"
        />
        <SunIcon class="swap-on size-6" />
        <MoonIcon class="swap-off size-6" />
      </label>
      <button class="logout-btn btn btn-outline btn-error" @click="onLogout">
        <ArrowRightStartOnRectangleIcon class="size-5" />
        {{ $t("account.logout") }}
      </button>
    </div>

    <form @submit.prevent="onSubmit" class="flex flex-1 flex-col" aria-label="Profil">
      <div class="flex w-full max-w-sm flex-col gap-6">
        <div class="flex flex-col gap-2">
          <label for="profile-username" class="font-medium">{{ $t("account.username") }}</label>
          <label class="input">
            <svg class="h-[1em] opacity-50" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
              <g
                stroke-linejoin="round"
                stroke-linecap="round"
                stroke-width="2.5"
                fill="none"
                stroke="currentColor"
              >
                <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"></path>
                <circle cx="12" cy="7" r="4"></circle>
              </g>
            </svg>
            <input
              id="profile-username"
              v-model="username"
              v-bind="usernameAttrs"
              type="text"
              :placeholder="$t('account.usernamePlaceholder')"
              autocomplete="username"
            />
            <span v-if="errors.username" class="text-sm text-error">{{ errors.username }}</span>
          </label>
        </div>

        <div class="flex flex-col gap-2">
          <label for="profile-password" class="font-medium">{{ $t("account.password") }}</label>
          <label class="input">
            <svg class="h-[1em] opacity-50" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
              <g
                stroke-linejoin="round"
                stroke-linecap="round"
                stroke-width="2.5"
                fill="none"
                stroke="currentColor"
              >
                <path
                  d="M2.586 17.414A2 2 0 0 0 2 18.828V21a1 1 0 0 0 1 1h3a1 1 0 0 0 1-1v-1a1 1 0 0 1 1-1h1a1 1 0 0 0 1-1v-1a1 1 0 0 1 1-1h.172a2 2 0 0 0 1.414-.586l.814-.814a6.5 6.5 0 1 0-4-4z"
                ></path>
                <circle cx="16.5" cy="7.5" r=".5" fill="currentColor"></circle>
              </g>
            </svg>
            <input
              id="profile-password"
              v-model="password"
              v-bind="passwordAttrs"
              type="password"
              :placeholder="$t('account.passwordPlaceholder')"
              autocomplete="new-password"
            />
          </label>
          <span v-if="errors.password" class="text-sm text-error">{{ errors.password }}</span>
        </div>

        <div class="flex flex-col gap-2">
          <label for="profile-password-confirm" class="font-medium">{{ $t("account.confirmPassword") }}</label>
          <label class="input">
            <svg class="h-[1em] opacity-50" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
              <g
                stroke-linejoin="round"
                stroke-linecap="round"
                stroke-width="2.5"
                fill="none"
                stroke="currentColor"
              >
                <path
                  d="M2.586 17.414A2 2 0 0 0 2 18.828V21a1 1 0 0 0 1 1h3a1 1 0 0 0 1-1v-1a1 1 0 0 1 1-1h1a1 1 0 0 0 1-1v-1a1 1 0 0 1 1-1h.172a2 2 0 0 0 1.414-.586l.814-.814a6.5 6.5 0 1 0-4-4z"
                ></path>
                <circle cx="16.5" cy="7.5" r=".5" fill="currentColor"></circle>
              </g>
            </svg>
            <input
              id="profile-password-confirm"
              v-model="confirmPassword"
              v-bind="confirmPasswordAttrs"
              type="password"
              :placeholder="$t('account.confirmPasswordPlaceholder')"
              autocomplete="new-password"
            />
          </label>
          <span v-if="errors.confirmPassword" class="text-sm text-error">{{ errors.confirmPassword }}</span>
        </div>
      </div>

      <div class="mt-auto flex justify-end">
        <button type="submit" :disabled="loading || !meta.dirty || !meta.valid" class="btn btn-success">
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
