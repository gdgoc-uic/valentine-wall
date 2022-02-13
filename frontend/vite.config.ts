import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
// import WindiCSS from 'vite-plugin-windicss'
import Icons from 'unplugin-icons/vite'
import { FileSystemIconLoader } from 'unplugin-icons/loaders'
import checker from 'vite-plugin-checker'
import { isoImport } from 'vite-plugin-iso-import'
import markdown, { Mode } from 'vite-plugin-markdown'
import { ViteFaviconsPlugin } from 'vite-plugin-favicon2'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    Icons({
      customCollections: {
        'home-icons': FileSystemIconLoader(
          './src/assets/images/home/icons',
        ),
      }
    }),
    checker({typescript: true}),
    isoImport(),
    markdown({
      mode: [Mode.HTML, Mode.VUE]
    }),
    ...(process.env.NODE_ENV == 'production' ? [
      ViteFaviconsPlugin({
        logo: "src/assets/images/icon.png"
      })
    ] : [])
  ]
})
