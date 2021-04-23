<!--
     This program is free software: you can redistribute it and/or modify
     it under the terms of the Affero GNU General Public License as published by
     the Free Software Foundation, either version 3 of the License, or
     (at your option) any later version.

     This program is distributed in the hope that it will be useful,
     but WITHOUT ANY WARRANTY; without even the implied warranty of
     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
     GNU General Public License for more details.

     You should have received a copy of the Affero GNU General Public License
     along with this program.  If not, see <https://www.gnu.org/licenses/>.
-->

<template>
    <div>
      <headerDisplay/>
      <div class="container">
        <section class="section form-width">
          <h1 class="title is-1">Login</h1>
          <p class="subtitle is-4">
              Welcome to FlatTrack, please login.
          </p>
          <b-field label="Email" class="is-marginless">
            <b-input
              name="email"
              type="email"
              v-model="email"
              maxlength="70"
              autofocus
              placeholder="Enter your email"
              @keyup.enter.native="postLogin"
              size="is-medium"
              icon="email"
              icon-right="close-circle"
              icon-right-clickable
              @icon-right-click="email = ''"
              required>
              </b-input>
          </b-field>
          <b-field label="Password" class="is-marginless">
            <b-input
              name="password"
              type="password"
              v-model="password"
              password-reveal
              maxlength="70"
              @keyup.enter.native="postLogin"
              placeholder="Enter your password"
              size="is-medium"
              icon="textbox-password"
              pattern="^([a-z]*)([A-Z]*).{10,}$"
              icon-right="close-circle"
              icon-right-clickable
              @icon-right-click="password = ''"
              required>
            </b-input>
          </b-field>
          <div class="field">
            <p class="control">
              <b-button
                icon-left="login"
                native-type="submit"
                size="is-medium"
                type="is-primary"
                expanded
                @click="postLogin">
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
            <div class="notification is-warning mb-4 mt-2" v-if="typeof message !== 'undefined' && message !== ''">
              <p class="subtitle is-6">{{  message }}</p>
            </div>
          </div>
        </section>
      </div>
    </div>
</template>

<script>
import { LoadingProgrammatic as Loading } from 'buefy'
import login from '@/requests/public/login'
import common from '@/common/common'

export default {
  name: 'login',
  data () {
    return {
      redirect: this.$route.query.redirect || null,
      authToken: this.$route.query.authToken || null,
      message: common.GetLoginMessage(),
      email: '',
      password: ''
    }
  },
  components: {
    headerDisplay: () => import('@/components/common/header-display')
  },
  methods: {
    postLogin () {
      const loadingComponent = Loading.open({
        container: null
      })
      setTimeout(() => loadingComponent.close(), 20 * 1000)
      login.PostUserAuth(this.email, this.password)
        .then(resp => {
          common.SetAuthToken(resp.data.data)
          setTimeout(() => {
            loadingComponent.close()
            if (this.redirect !== null) {
              this.$router.push({ path: this.redirect })
              return
            }
            window.location.href = '/'
          }, 2 * 1000)
        }).catch(err => {
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

            if (this.redirect !== null) {
              this.$router.push({ path: this.redirect })
              return
            }
            window.location.href = '/'
          }, 2 * 1000)
        })
      }
    }
  },
  mounted () {
    if (typeof this.authToken !== 'undefined') {
      common.SetAuthToken(this.authToken)
    }
    this.checkForLoginToken()
  }
}
</script>

<style scoped>
</style>
