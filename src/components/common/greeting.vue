<template>
    <div>
        <h2 class="title is-2">Hey, {{ login.names }}!</h2>
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
  name: 'greeting',
  data () {
    return {
      login: {
        names: ''
      }
    }
  },
  created () {
    axios.get(`/api/profile`,
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('authToken')}`
        }
      })
      .then(resp => {
        this.login = resp.data
      })
      .catch(err => {
        this.$buefy.notification.open({
          duration: 5000,
          message: `An error has occured: ${err}`,
          position: 'is-bottom-right',
          type: 'is-danger',
          hasIcon: true
        })
      })
  }
}
</script>

<style scoped>

</style>
