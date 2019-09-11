<template>
  <div>
    <section class="hero is-info">
      <div class="hero-body">
        <p class="title">
          FlatTrack
        </p>
        <p class="subtitle">
          {{ deploymentName }}
        </p>
      </div>
    </section>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'header-display',
  data () {
    return {
      deploymentName: 'Keep track of your flat'
    }
  },
  created () {
    axios.get(`/api/settings/deploymentName`)
      .then(response => {
        this.deploymentName = response.data.value
      })
      .catch(err => {
        this.pageErrors.push(err)
      })
  }
}
</script>

<style>

</style>