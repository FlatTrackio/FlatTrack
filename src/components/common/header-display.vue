<template>
  <div>
    <section :class="`${admin !== 'true' ? 'hero is-info' : 'hero is-warning'}`">
      <div class="hero-body">
        <p class="title">
          FlatTrack {{ admin !== 'true' ? '' : '(Admin)' }}
        </p>
        <p class="subtitle">
          <span v-if="genericMessage === 'true' || admin === 'true'">
            Keep track of your flat
          </span>
          <span v-else>
            {{ subtitle }}
          </span>
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
    'genericMessage': Boolean,
    'admin': Boolean
  },
  data () {
    return {
      subtitle: 'Keep track of your flat'
    }
  },
  created () {
    if (this.genericMessage !== 'true') {
      axios.get(`/api/settings/deploymentName`)
        .then(response => {
          this.subtitle = response.data.value
        })
        .catch(err => {
          this.pageErrors.push(err)
        })
    }
  }
}
</script>

<style scoped>

</style>
