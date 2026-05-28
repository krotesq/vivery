<script setup lang="ts">
import { RouterLink } from 'vue-router';
</script>

<template>
<div class="card bg-base-100 shadow-sm w-full md:w-96">
    <div class="card-body">
        <p>{{ $t("notFound.title") }}</p>
        <div class="card-actions">
            <RouterLink to="/" class="btn btn-primary w-full">
            {{ $t("notFound.backHome") }}
            </RouterLink>
        </div>
    </div>
</div>
</template>