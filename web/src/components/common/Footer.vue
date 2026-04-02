<template>
  <footer class="bg-gradient-to-r from-gray-900 via-gray-800 to-gray-900 text-white">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <!-- 关于 -->
        <div>
          <h3 class="text-xl font-bold mb-4 flex items-center gap-2">
            <span class="text-2xl">🌊</span>
            Cyblog
          </h3>
          <p class="text-gray-400 text-sm leading-relaxed">
            一个简洁优雅的博客系统，使用 Go 和 Vue3 构建。分享技术，记录生活。
          </p>
        </div>

        <!-- 快速链接 -->
        <div>
          <h3 class="text-lg font-bold mb-4">快速链接</h3>
          <ul class="space-y-2">
            <li>
              <router-link to="/" class="text-gray-400 hover:text-white transition-colors text-sm">
                首页
              </router-link>
            </li>
            <li>
              <router-link to="/article" class="text-gray-400 hover:text-white transition-colors text-sm">
                文章列表
              </router-link>
            </li>
            <li>
              <router-link to="/about" class="text-gray-400 hover:text-white transition-colors text-sm">
                关于我们
              </router-link>
            </li>
          </ul>
        </div>

        <!-- 联系 -->
        <div>
          <h3 class="text-lg font-bold mb-4">联系我们</h3>
          <ul class="space-y-2 text-sm text-gray-400">
            <li class="flex items-center gap-2">
              <span>📧</span>
              <span>contact@cyblog.com</span>
            </li>
            <li class="flex items-center gap-2">
              <span>🌐</span>
              <span>www.cyblog.com</span>
            </li>
          </ul>
          <div class="flex gap-4 mt-4">
            <a href="#" class="w-10 h-10 bg-white/10 rounded-full flex items-center justify-center hover:bg-white/20 transition-colors">
              <span>📱</span>
            </a>
            <a href="#" class="w-10 h-10 bg-white/10 rounded-full flex items-center justify-center hover:bg-white/20 transition-colors">
              <span>📧</span>
            </a>
            <a href="#" class="w-10 h-10 bg-white/10 rounded-full flex items-center justify-center hover:bg-white/20 transition-colors">
              <span>🐙</span>
            </a>
          </div>
        </div>
      </div>

      <div class="border-t border-gray-700 mt-8 pt-8 text-center text-gray-500 text-sm">
        <p>© 2024 Cyblog. 使用 ❤️ 和 ☕ 构建</p>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
defineOptions({
  name: 'Footer',
})
</script>
