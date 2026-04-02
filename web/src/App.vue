<script setup lang="ts">
import ContentLayout from './components/layout/ContentLayout.vue'
</script>

<template>
  <ContentLayout>
    <RouterView />
  </ContentLayout>
</template>

<style>
#app {
  min-height: 100vh;
}
</style>
