<template>
  <footer class="relative bg-gradient-to-r from-gray-900 via-sky-900 to-gray-900 text-white py-16 mt-20">
    <div class="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjAiIGhlaWdodD0iNjAiIHZpZXdCb3g9IjAgMCA2MCA2MCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48ZyBmaWxsPSJub25lIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiPjxnIGZpbGw9IiNmZmYiIGZpbGwtb3BhY2l0eT0iMC4wNSI+PGNpcmNsZSBjeD0iMzAiIGN5PSIzMCIgcj0iMiIvPjwvZz48L2c+PC9zdmc+')] opacity-50"></div>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
      <div class="text-center">
        <div class="text-3xl font-bold bg-gradient-to-r from-sky-300 to-cyan-300 bg-clip-text text-transparent mb-4">
          🌊 Cyblog
        </div>
        <p class="text-gray-400 mb-6">
          用心分享，每一篇文章都是一份礼物 🎁
        </p>
        <div class="flex justify-center gap-4 mb-8">
          <router-link to="/" class="w-10 h-10 bg-white/10 rounded-full flex items-center justify-center hover:bg-white/20 transition-colors">
            🏠
          </router-link>
          <router-link to="/article" class="w-10 h-10 bg-white/10 rounded-full flex items-center justify-center hover:bg-white/20 transition-colors">
            📝
          </router-link>
          <router-link to="/category" class="w-10 h-10 bg-white/10 rounded-full flex items-center justify-center hover:bg-white/20 transition-colors">
            🏷️
          </router-link>
        </div>
        <div class="border-t border-white/10 pt-8">
          <p class="text-gray-500">
            © 2026 Cyblog. Made with 💖 and ☕
          </p>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
// Footer 组件无需额外逻辑
</script>
