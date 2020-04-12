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
          <!-- TODO add FlatTrack version and links -->
        </section>
      </div>
    </div>
</template>

<script>
import { LoadingProgrammatic as Loading } from 'buefy'
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
    headerDisplay: () => import('@/frontend/components/header-display')
  },
  methods: {
    postLogin (email, password) {
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
            window.location.href = '/'
          }, 2 * 1000)
        }).catch(err => {
          console.log(err)
          loadingComponent.close()
          common.DisplayFailureToast(err.response.data.metadata.response || err)
        })
    },
    checkForLoginToken () {
      var authToken = common.GetAuthToken()
      if (!(typeof authToken === 'undefined' || authToken === null || authToken === '')) {
        login.GetUserAuth(false).then(res => {
          // verify token via request or something
          const loadingComponent = Loading.open({
            container: null
          })
          common.DisplaySuccessToast('You are still signed in, going to the home page...')
          setTimeout(() => {
            loadingComponent.close()
            window.location.href = '/'
          }, 2 * 1000)
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
</style>
