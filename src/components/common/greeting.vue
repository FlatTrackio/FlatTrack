<template>
    <div>
        <h2 class="title is-2">Hey, {{ login.names }}!</h2>
    </div>
</template>

<script>
import axios from 'axios'

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
      .then(response => {
        this.login = response.data
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
