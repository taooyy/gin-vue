import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
  ],
  resolve: {
    alias: {
      // 这里配置了 @ 指向 src 目录，解决 "找不到模块" 的核心配置
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
  // 我把 css 配置块删掉了，这样就不会报错了
})