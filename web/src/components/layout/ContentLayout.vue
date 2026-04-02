<template>
  <div class="content-layout">
    <Header />
    <main class="content-main">
      <RouterView></RouterView>
    </main>
    <Footer />
  </div>
</template>

<script setup lang="ts">
import Header from '@/components/common/Header.vue'
import Footer from '@/components/common/Footer.vue'

defineOptions({
  name: 'ContentLayout',
})
</script>

<style scoped>
.content-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.content-main {
  flex: 1;
}
</style>
