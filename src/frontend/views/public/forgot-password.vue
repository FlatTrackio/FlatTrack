<template>
    <div>
      <headerDisplay/>
      <div class="container">
        <section class="section form-width">
          <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
          <ul>
              <li><router-link to="/login">Login</router-link></li>
              <li class="is-active"><router-link to="/forgot-password">Forgot Password</router-link></li>
          </ul>
          </nav>
          <h1 class="title is-1">Forgot Password</h1>
          <p class="subtitle is-4">Enter you email to reset your password</p>
          <b-field label="Email" class="is-marginless">
            <b-input
              type="email"
              v-model="email"
              maxlength="70"
              autofocus
              placeholder="Enter your email"
              @keyup.enter.native="sendPasswordResetRequest(email)"
              size="is-medium"
              icon="email"
              required>
              </b-input>
          </b-field>
          <b-button
            expanded
            @click="sendPasswordResetRequest(email)"
            size="is-medium"
            type="is-primary">
            Reset
          </b-button>
        </section>
      </div>
    </div>
</template>

<script>
import headerDisplay from '@/frontend/components/common/header-display'
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
      // post to api to request a password reset email
    }
  },
  components: {
    headerDisplay
  }
}
</script>

<style scoped>

</style>
