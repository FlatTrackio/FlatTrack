<template>
  <div>
    <section class="hero is-info">
      <div class="hero-body">
        <p class="title">
          FlatTrack
        </p>
        <p class="subtitle">
          <div v-if="genericMessage == 'true'">
            Keep track of your flat
          </div>
          <div v-else>
            {{ subtitle }}
          </div>
        </p>
      </div>
    </section>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'header-display',
  props: {
    'genericMessage': Boolean
  },
  data () {
    return {
      subtitle: 'Keep track of your flat'
    }
  },
  created () {
    axios.get(`/api/settings/deploymentName`)
      .then(response => {
        this.subtitle = response.data.value
      })
      .catch(err => {
        this.pageErrors.push(err)
      })
  }
}
</script>

<style>

</style>
