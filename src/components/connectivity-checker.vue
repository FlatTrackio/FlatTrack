<template>
    <div>
    </div>
</template>

<script>
import axios from 'axios'
import { Service } from 'axios-middleware'
import { NotificationProgrammatic as Notification } from 'buefy'

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
  name: 'connectivity-checker',
  data () {
    return {
      active: true
    }
  },
  created () {
    function sleep (ms) {
      return new Promise(resolve => setTimeout(resolve, ms))
    }

    (async () => {
      while (true) {
        await axios.get('/api').then(() => {
          if (this.active === false) {
            this.active = true
            Notification.open({
              duration: 5000,
              message: `Hooray! I'm talking to the server again.`,
              position: 'is-bottom-right',
              type: 'is-success',
              hasIcon: true
            })
          }
        }).catch(err => {
          this.active = false
          console.log({err})
          Notification.open({
            duration: 10000,
            message: `Uh oh, I'm having an issue talking to the server.`,
            position: 'is-bottom-right',
            type: 'is-danger',
            hasIcon: true
          })
        })
        await sleep(10 * 1000)
      }
    })()
  }
}
</script>

<style scoped>

</style>
