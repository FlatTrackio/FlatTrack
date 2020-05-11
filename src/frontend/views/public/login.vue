<template>
    <div>
      <headerDisplay/>
      <div class="container">
        <section class="section form-width">
          <h1 class="title is-1">Login</h1>
          <p class="subtitle is-3">
              Welcome to FlatTrack, please login.
          </p>
          <b-field label="Email" class="is-marginless">
            <b-input
              type="email"
              v-model="email"
              maxlength="70"
              autofocus
              placeholder="Enter your email"
              @keyup.enter.native="postLogin(email, password)"
              size="is-medium"
              icon="email"
              required>
              </b-input>
          </b-field>
          <b-field label="Password">
            <b-input
              type="password"
              v-model="password"
              password-reveal
              maxlength="70"
              @keyup.enter.native="postLogin(email, password)"
              placeholder="Enter your password"
              size="is-medium"
              icon="textbox-password"
              pattern="^([a-z]*)([A-Z]*).{10,}$"
              validation-message="Password is invalid. Passwords must include: one number, one lowercase letter, one uppercase letter, and be eight or more characters."
              required>
            </b-input>
          </b-field>
          <div class="field">
            <p class="control">
              <b-button
                icon-left="login"
                native-type="submit"
                size="is-medium"
                expanded
                @click="postLogin(email, password)">
                Login
              </b-button>
              <b-button
                tag="a"
                href="forgot-password"
                icon-left="lifebuoy"
                size="is-medium"
                expanded
                disabled
                type="is-text">
                Forgot Password
              </b-button>
            </p>
          </div>
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
    headerDisplay: () => import('@/frontend/components/common/header-display')
  },
  methods: {
    postLogin (email, password) {
      console.log({ email, password })
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
