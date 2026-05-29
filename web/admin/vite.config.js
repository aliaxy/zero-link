import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const apiTarget = env.VITE_API_TARGET || 'http://127.0.0.1:8080'

  return {
    plugins: [
      vue(),
      Components({
        resolvers: [ElementPlusResolver()],
      }),
    ],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
    server: {
      port: 5173,
      proxy: {
        // Admin API calls: /api/* → link-api (strip /api prefix)
        '/api': {
          target: apiTarget,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api/, ''),
        },
        // Short link redirects: /:code → link-api
        // Excludes known SPA routes: login, links, api
        // $ anchor prevents matching /src/main.js and other asset paths
        '^/(?!api(/|$)|login$|links$)[a-zA-Z0-9_-]{3,32}$': {
          target: apiTarget,
          changeOrigin: true,
        },
      },
    },
  }
})
