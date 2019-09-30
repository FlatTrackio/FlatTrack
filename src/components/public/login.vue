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
                  v-model="email"
                  maxlength="70"
                  class="is-focused"
                  required>
              </b-input>
          </b-field>
          <b-field label="Password">
              <b-input type="password"
                  v-model="password"
                  password-reveal
                  maxlength="70"
                  required>
              </b-input>
          </b-field>
          <b-button rounded native-type="submit" @click="postLogin(email, password)">Login</b-button>
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
      email: '',
      password: ''
    }
  },
  components: {
    headerDisplay
  },
  methods: {
    postLogin: (email, password) => {
      const loadingComponent = Loading.open({
        container: null
      })
      setTimeout(() => loadingComponent.close(), 20 * 1000)
      var form = {
        email: email,
        password: password
      }
      axios.post('/api/login', form,
        {
          headers: {
            Authorization: `Bearer ${sessionStorage.getItem('authToken')}`
          }
        }).then(resp => {
        console.log('Signing in')
        sessionStorage.setItem('authToken', resp.data.refreshToken)
        setTimeout(() => {
          loadingComponent.close()
          location.href = '/'
        }, 2000)
      }).catch(err => {
        console.log('Failed to sign in')
        console.log(err)
        loadingComponent.close()
        Toast.open({
          duration: 10000,
          message: 'Hmmm, something went wrong with the login. Email or password is incorrect. Please try again.',
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
          loadingComponent.close()
          console.log('Found auth token')
          location.href = '/'
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
