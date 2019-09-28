<template>
    <div>
        <headerDisplay genericMessage="true"/>
        <div class="container">
          <section class="section">
              <h1 class="title">Login</h1>
              <p class="subtitle">
                  Welcome to FlatTrack, please login.
              </p>
              <b-field label="Email">
                  <b-input type="email"
                      :value="form.email"
                      maxlength="70"
                      class="is-focused"
                      required>
                  </b-input>
              </b-field>
              <b-field label="Password">
                  <b-input type="password"
                      :value="form.password"
                      password-reveal
                      maxlength="70"
                      required>
                  </b-input>
              </b-field>
              <b-button rounded @click="postLogin(form)">Login</b-button>
              <b-button tag="a" href="#/forgot-password" rounded type="is-warning">Forgot Password</b-button>
          </section>
        </div>
    </div>
</template>

<script>
import axios from 'axios'
import headerDisplay from '../common/header-display'
import { ToastProgrammatic as Toast, LoadingProgrammatic as Loading } from 'buefy'

export default {
  name: 'login',
  data () {
    return {
      form: {
        email: '',
        password: ''
      }
    }
  },
  components: {
    headerDisplay
  },
  methods: {
    postLogin: (form) => {
      const loadingComponent = Loading.open({
        container: null
      })
      setTimeout(() => loadingComponent.close(), 50 * 1000)
      form = {
        email: form.email,
        password: form.password
      }
      console.log(JSON.stringify(form))
      axios.post('/api/login', form).then(resp => {
        console.log(resp)
      }).catch(err => {
        console.log(err)
        loadingComponent.close()
        Toast.open({
          duration: 10000,
          message: 'Hmmm, something went wrong with the login. Email or password is incorrect. Please type again',
          position: 'is-bottom',
          type: 'is-danger'
        })
      })
    },
    checkForLoginToken: () => {
      if (sessionStorage.getItem('authToken')) {
        // verify token via request or something
        const loadingComponent = Loading.open({
          container: null
        })
        Toast.open({
          duration: 2000,
          message: 'You\'re still signed in, let\'s go to the home page',
          position: 'is-bottom'
        })
        setTimeout(() => {
          console.log('Found auth token')
          location.href = '/'
          loadingComponent.close()
        }, 2000)
      }
    }
  },
  mounted () {
    this.checkForLoginToken()
  }
}
</script>

<style scoped>

</style>
