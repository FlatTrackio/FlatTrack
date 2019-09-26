<template>
    <div>
        <headerDisplay genericMessage="true"/>
        <section class="section">
            <nav class="breadcrumb has-arrow-separator" aria-label="breadcrumbs">
            <ul>
                <li><a href="/#/login">Login</a></li>
                <li class="is-active"><a href="/#/forgot-password">Forgot Password</a></li>
            </ul>
            </nav>
            <h1 class="title">Forgot Password</h1>
            <h2 class="subtitle">Enter you email to reset your password</h2>
        </section>
        <section class="section">
            <b-field label="Email">
                <b-input type="email"
                    :value="email"
                    maxlength="70">
                </b-input>
            </b-field>
            <b-button rounded @click="sendPasswordResetRequest(email)" type="is-light">Reset</b-button>
        </section>
    </div>
</template>

<script>
import axios from 'axios'
import headerDisplay from '../common/header-display'
import { ToastProgrammatic as Toast } from 'buefy'

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
      axios.post('/api/')
    }
  },
  components: {
    headerDisplay
  }
}
</script>

<style scoped>

</style>
