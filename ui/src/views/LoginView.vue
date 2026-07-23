<script setup lang="ts">
  import { useRouter, useRoute } from 'vue-router'
  import { useAuthStore } from '@/stores/auth'
  import { useForm } from 'vee-validate'
  import { toTypedSchema } from '@vee-validate/zod'
  import { loginSchema } from '@/schemas/account'
  import { ref } from 'vue'
  import { toast } from 'vue-sonner'

  const router = useRouter()
  const route = useRoute()
  const auth = useAuthStore()

  const { handleSubmit, defineField } = useForm({
    validationSchema: toTypedSchema(loginSchema),
  })

  const [username, usernameAttrs] = defineField('username')
  const [password, passwordAttrs] = defineField('password')

  const loading = ref(false)

  const onSubmit = handleSubmit(async (values) => {
    loading.value = true
    try {
      await auth.login(values)
      const redirect = route.query.redirect as string | undefined
      await router.push(redirect ?? '/')
    } catch (e) {
      toast.error(e instanceof Error ? e.message : 'Error while logging in')
    } finally {
      loading.value = false
    }
  })
</script>

<template>
  <main class="flex min-h-screen w-screen items-center justify-center bg-base-100 p-4 text-base-content sm:p-8 lg:items-stretch lg:justify-between lg:gap-10 lg:p-15">
    <section class="flex min-h-[calc(100vh-2rem)] w-full flex-col items-center justify-center gap-12 rounded-box bg-base-200 p-6 sm:min-h-[calc(100vh-4rem)] sm:p-10 lg:min-h-0 lg:w-100 lg:gap-30 lg:p-15 shadow-[15px_15px_rgba(0,0,0,0.25)]">
      <header class="w-full text-center">
        <h1 class="text-2xl font-bold text-shadow-[5px_5px_rgb(0_0_0/0.25)]">Vivery v0.0.1-alpha</h1>
      </header>

      <form @submit.prevent="onSubmit" class="flex w-full max-w-80 flex-col gap-4 sm:w-80 lg:w-60" aria-label="Login">
        <label class="sr-only" for="login-username">{{ $t('account.username') }}</label>
        <input
          id="login-username"
          v-model="username"
          v-bind="usernameAttrs"
          type="text"
          :placeholder="$t('account.username')"
          autocomplete="username"
          class="input input-primary w-full shadow-[5px_5px_rgba(0,0,0,0.25)]"
        />

        <label class="sr-only" for="login-password">{{ $t('account.password') }}</label>
        <input
          id="login-password"
          v-model="password"
          v-bind="passwordAttrs"
          type="password"
          :placeholder="$t('account.password')"
          autocomplete="current-password"
          class="input input-primary w-full shadow-[5px_5px_rgba(0,0,0,0.25)]"
        />

        <button type="submit" :disabled="loading" class="btn btn-soft btn-primary w-full shadow-[5px_5px_rgba(0,0,0,0.25)]">
          {{ loading ? `${$t("misc.loading")}...` : $t("account.login") }}
        </button>
      </form>
    </section>

    <section class="hidden min-h-80 flex-1 rounded-box bg-base-100 lg:flex lg:items-center lg:justify-center" aria-label="Animation">
      Animation
    </section>
  </main>
</template>
