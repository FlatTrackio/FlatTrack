import Vue from 'vue'
import App from './App.vue'
import router from './router'
import Buefy from 'buefy'
import VueMaterial from 'vue-material'
import 'vue-material/dist/vue-material.min.css'
import './registerServiceWorker'

Vue.use(VueMaterial)
Vue.use(Buefy, {
  defaultIconPack: 'mdi'
})
Vue.config.productionTip = false

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
