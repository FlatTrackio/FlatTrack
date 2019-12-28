<template>
  <div>
    <section :class="`${admin !== 'true' ? 'hero is-info' : 'hero is-warning'}`">
      <div class="hero-body">
        <p class="title">
          FlatTrack {{ admin !== 'true' ? '' : '(Admin)' }}
        </p>
        <p class="subtitle">
          {{ subtitle }}
        </p>
      </div>
    </section>
  </div>
</template>

<script>
import axios from 'axios'
import { Service } from 'axios-middleware'

const service = new Service(axios)
service.register({
  onResponse (response) {
    if (response.status === 403) {
      localStorage.removeItem('authToken')
      location.href = '/'
    }
    return response
  }
})

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
      axios.get(`/api/settings/deploymentName`,
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('authToken')}`
          }
        }).then(resp => {
        this.subtitle = resp.data.value
      }).catch(err => {
        this.pageErrors.push(err)
      })
    }
  }
}
</script>

<style scoped>

</style>
