<template>
  <div>
    <headerDisplay/>
    <div class="container">
      <section class="section">
        <h1 class="title is-1">Confirm your account</h1>
        <p class="subtitle is-3">
          Final things to complete your sign up
        </p>

        <b-message type="is-danger" has-icon v-if="(typeof id === 'undefined' || id === '') || idValid !== true">
          Token Id is missing or is invalid
          <!-- TODO explain better -->
        </b-message>
        <b-message type="is-danger" has-icon v-if="typeof secret === 'undefined' || secret === ''">
          Missing a confirmation secret.
          <!-- TODO explain better -->
        </b-message>
        <div v-if="idValid === true && typeof secret !== 'undefined' && secret !== ''">
          <b-field label="Phone number*">
            <b-input
              type="tel"
              v-model="phoneNumber"
              placeholder="Enter your phone number"
              icon="phone"
              size="is-medium"
              maxlength="30">
            </b-input>
          </b-field>

          <b-field label="Birthday*">
            <b-datepicker
              v-model="jsBirthday"
              :max-date="maxDate"
              :show-week-numbers="true"
              :focused-date="focusedDate"
              placeholder="Click to select birthday"
              icon="cake-variant"
              size="is-medium"
              trap-focus>
            </b-datepicker>
          </b-field>
          <br/>

          <b-field label="Password">
            <b-input
              type="password"
              v-model="password"
              password-reveal
              maxlength="70"
              placeholder="Enter your password"
              icon="textbox-password"
              size="is-medium"
              required>
            </b-input>
          </b-field>

          <b-field label="Confirm password">
            <b-input
              type="password"
              v-model="passwordConfirm"
              password-reveal
              maxlength="70"
              placeholder="Confirm your password"
              icon="textbox-password"
              size="is-medium"
              required>
            </b-input>
          </b-field>
          <b-button
            type="is-success"
            size="is-medium"
            icon-left="check"
            native-type="submit"
            @click="PostUserConfirm(id, secret, phoneNumber, password, passwordConfirm, jsBirthday)">
            Confirm my account
          </b-button>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import { LoadingProgrammatic as Loading } from 'buefy'
import common from '@/frontend/common/common'
import confirm from '@/frontend/requests/public/useraccountconfirm'
import moment from 'moment'

export default {
  name: 'Account confirmation',
  data () {
    var today = new Date()
    var maxDate = new Date(today.getFullYear() - 15, today.getMonth(), today.getDay())

    return {
      idValid: false,
      jsBirthday: null,
      id: this.$route.params.id,
      secret: this.$route.query.secret,
      phoneNumber: null,
      password: null,
      passwordConfirm: null
    }
  },
  components: {
    headerDisplay: () => import('@/frontend/components/common/header-display')
  },
  methods: {
    PostUserConfirm (id, secret, phoneNumber, password, passwordConfirm, jsBirthday) {
      if (password !== passwordConfirm && password !== '') {
        common.DisplayFailureToast('Passwords do not match')
        return
      }
      var birthday = Number(moment(jsBirthday).format('X')) || 0

      const loadingComponent = Loading.open({
        container: null
      })
      confirm.PostUserConfirm(id, secret, phoneNumber, birthday, password).then(resp => {
        if (resp.data.data === '') {
          common.DisplayFailureToast('Unable to find token to sign in with after confirming account. Please contact an administrator')
          return
        }
        localStorage.setItem('authToken', resp.data.data)
        setTimeout(() => {
          loadingComponent.close()
          window.location.href = '/'
        }, 2 * 1000)
      }).catch(() => {
        loadingComponent.close()
        common.DisplayFailureToast('Unable to find token to sign in with after confirming account. Please contact an administrator')
      })
    }
  },
  async beforeMount () {
    confirm.GetTokenValid(this.id).then(resp => {
      this.idValid = resp.data.data
    })
  }
}
</script>

<style scoped>

</style>
