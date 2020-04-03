<template>
    <div>
      <headerDisplay/>
      <div class="container">
        <section class="section form-width">
          <h1 class="title">Login</h1>
          <p class="subtitle">
              Welcome to FlatTrack, please login.
          </p>
          <b-field label="Email">
              <b-input type="email"
                  v-model="email"
                  maxlength="70"
                  autofocus
                  required>
              </b-input>
          </b-field>
          <b-field label="Password">
              <b-input type="password"
                  v-model="password"
                  password-reveal
                  maxlength="70"
                  @keyup.enter.native="postLogin(email, password)"
                  required>
              </b-input>
          </b-field>
          <b-button rounded native-type="submit" @click="postLogin(email, password)">Login</b-button>
          <b-button tag="a" href="forgot-password" rounded type="is-warning">Forgot Password</b-button>
        </section>
      </div>
    </div>
</template>

<script>
import headerDisplay from '@/frontend/components/header-display'
import { ToastProgrammatic as Toast, LoadingProgrammatic as Loading } from 'buefy'
import login from '@/frontend/requests/public/login'
import common from '@/frontend/common/common'

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
      login.PostUserAuth(form)
        .then(resp => {
          localStorage.setItem('authToken', resp.data.data)
          setTimeout(() => {
            loadingComponent.close()
            location.href = '/'
          }, 2000)
        }).catch(err => {
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
      var authToken = common.GetAuthToken()
      if (!(typeof authToken === 'undefined' || authToken === null || authToken === '')) {
        login.GetUserAuth(false).then(res => {
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
            location.href = '/'
          }, 2000)
        })
      }
    }
  },
  mounted () {
    this.checkForLoginToken()
  }
}
</script>

<style scoped>
.form-width {
    width: 380px;
    margin: auto;
}
</style>
