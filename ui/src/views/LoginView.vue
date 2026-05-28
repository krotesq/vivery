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

  <form @submit.prevent="onSubmit" class="card bg-base-100 shadow-sm w-full md:w-96">

    <div class="card-body">

      <h2 class="text-xl">{{ $t("account.login") }}</h2>

        <input v-model="username" v-bind="usernameAttrs" type="text" :placeholder="$t('account.username')" class="input w-full" />
        <input v-model="password" v-bind="passwordAttrs" type="password" :placeholder="$t('account.password')" class="input w-full" />
      
      <div class="card-actions">
        <button  type="submit" :disabled="loading" class="btn btn-primary w-full">{{ loading ? `${$t("misc.loading")}...` : $t("account.login") }}</button>
      </div>
    </div>
  </form>
</template>