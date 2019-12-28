<template>
    <div>
      <headerDisplay genericMessage="true"/>
      <div class="container">
        <section class="section">
          <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
          <ul>
              <li><a href="/login">Login</a></li>
              <li class="is-active"><a href="/forgot-password">Forgot Password</a></li>
          </ul>
          </nav>
          <h1 class="title">Forgot Password</h1>
          <h2 class="subtitle">Enter you email to reset your password</h2>
          <b-field label="Email">
              <b-input type="email"
                  v-model="email"
                  maxlength="70"
                  required>
              </b-input>
          </b-field>
          <b-button rounded @click="sendPasswordResetRequest(email)" type="is-light">Reset</b-button>
        </section>
      </div>
    </div>
</template>

<script>
import axios from 'axios'
import { Service } from 'axios-middleware'
import headerDisplay from '@/components/header-display'
import { ToastProgrammatic as Toast } from 'buefy'

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
  name: 'forgot-password',
  data () {
    return {
      email: ''
    }
  },
  methods: {
    sendPasswordResetRequest: (email) => {
      if (email === '') {
        Toast.open({
          message: 'Please provide a valid email address',
          position: 'is-bottom',
          type: 'is-danger'
        })
      }
      axios.post('/api/',
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('authToken')}`
          }
        })
    }
  },
  components: {
    headerDisplay
  }
}
</script>

<style scoped>

</style>
