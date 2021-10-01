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

import { defineConfig } from 'vite'
import { createVuePlugin } from 'vite-plugin-vue2'
const path = require('path')
export default defineConfig({
  plugins: [
    createVuePlugin(),
  ],
  server: {
    port: 8081
  },
  resolve: {
    alias: [
      {
        find: '@',
        replacement: path.resolve(__dirname, 'src')
      }
    ]
  },
  productionSourceMap: false,
  pwa: {
    name: "FlatTrack",
    themeColor: "#209cee",
    msTileColor: "#209cee",
    appleMobileWebAppCache: "yes",
    manifestOptions: {
      background_color: "#000000"
    },
    workboxOptions: {
      swSrc: 'service-worker.js',
      skipWaiting: true
    },
    workboxPluginMode: "GenerateSW"
  }
})
