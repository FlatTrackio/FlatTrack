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

import { createApp, getCurrentInstance } from "vue";
import { ToastProgrammatic as Toast } from "buefy";
import App from "./App.vue";
import router from "./router";
import Buefy from "buefy";
import "buefy/dist/css/buefy.css";
import { registerSW } from "virtual:pwa-register";

const intervalMS = 10 * 60 * 1000;

const updateSW = registerSW({
  onNeedRefresh() {
    new Toast(app).open({
      message: "FlatTrack update available. Please refresh",
      type: "is-info",
      position: "is-bottom",
    });
  },
  onRegisteredSW(swUrl, r) {
    r &&
      setInterval(async () => {
        if (r.installing || !navigator) return;
        if ("connection" in navigator && !navigator.onLine) return;
        const resp = await fetch(swUrl, {
          cache: "no-store",
          headers: {
            cache: "no-store",
            "cache-control": "no-cache",
          },
        });
        if (resp?.status === 200) await r.update();
      }, intervalMS);
  },
});

const app = createApp(App);
app.use(Buefy, {
  defaultIconPack: "mdi",
});
app.use(router);
app.mount("#app");
