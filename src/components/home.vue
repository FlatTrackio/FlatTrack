<template>
  <div>
  </div>
</template>

<script>
import VueRouter from 'vue-router'

const router = new VueRouter({
  mode: 'history',
  base: '/',
  routes: [
    {
      path: '/',
      component: () => {
        import('./authenticated/home')
      }
    }
  ]
})

export default {
  data () {
  },
  router
}
</script>

<style scoped>

</style>
