<template>
    <div>
      <headerDisplay/>
      <div class="container">
        <section class="section">
          <h1 class="title is-1">Connection error</h1>
          <b-message
            has-icon
            type="is-warning">
          <p class="subtitle is-4">
            Hmmm... something appears to have gone wrong with your FlatTrack instance.
            <br />
            <br />
            Please contact your system administrator.
            <br />
            <br />
          </p>
          </b-message>
            <b-button
              type="is-white"
              @click="Reload"
              size="is-medium"
              rounded
              expanded>
              Refresh
            </b-button>
        </section>
      </div>
    </div>
</template>

<script>
import common from '@/common/common'

export default {
  name: 'login',
  components: {
    headerDisplay: () => import('@/components/common/header-display')
  },
  methods: {
    Reload () {
      this.$router.go()
    }
  }
}
</script>

<style scoped>
</style>
