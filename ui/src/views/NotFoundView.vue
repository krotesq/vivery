<script setup lang="ts">
  import { RouterLink } from 'vue-router';

</script>

<template>

  <h1>{{ $t("notFound.title") }}</h1>
  <RouterLink to="/" class="btn btn-primary">
    {{ $t("notFound.backHome") }}
  </RouterLink>

</template>