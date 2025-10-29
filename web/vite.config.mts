// This program is free software: you can redistribute it and/or modify
// it under the terms of the Affero GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the Affero GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

import { defineConfig } from 'npm:vite@^5.2.10'
import vue from 'npm:@vitejs/plugin-vue@^5.0.4'

import 'npm:vue@3.2.26'
import "npm:vue-router@^4.6.3"
import "npm:axios@1.7.2"
import "npm:dayjs@1.11.11"
import 'npm:@mdi/js@7.4.47'
import 'npm:buefy@3.0.3'
import 'npm:canvas-confetti@1.4.0'
import 'npm:node-emoji@1.11.0'
import 'npm:qrcode.vue@1.7.0'
import 'npm:sass@1.93.2'
import 'npm:sass-loader@16.0.5'
import 'npm:register-service-worker@^1.7.2'
import 'npm:autoprefixer@^10.4.21'
import { NodePackageImporter } from 'npm:sass-embedded@1.93.2';
import { VitePWA } from "npm:vite-plugin-pwa@1.1.0";
import * as path from "jsr:@std/path"

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    VitePWA({
      registerType: "autoUpdate",
      includeAssets: [
        "favicon.svg",
        "favicon.ico",
        "robots.txt",
        "apple-touch-icon.png",
      ],
      manifest: {
        name: "FlatTrack",
        short_name: "FlatTrack",
        description: "Collaborate with your flatmates",
        theme_color: "#209cee",
        icons: [
          {
            src: "pwa-192x192.png",
            sizes: "192x192",
            type: "image/png",
          },
          {
            src: "pwa-512x512.png",
            sizes: "512x512",
            type: "image/png",
          },
          {
            src: "pwa-512x512.png",
            sizes: "512x512",
            type: "image/png",
            purpose: "any maskable",
          },
        ],
      },
    }),
  ],
  build: {
    assetsDir: "assets/",
  },
  server: {
    port: 8081,
  },
  resolve: {
    alias: [
      {
        find: "@",
        replacement: path.join(Deno.cwd(), "src"),
      },
    ],
  },
  css: {
    preprocessorOptions: {
      scss: {
        api: 'modern-compiler',
        importers: [new NodePackageImporter()],
      }
    }
  }
})
