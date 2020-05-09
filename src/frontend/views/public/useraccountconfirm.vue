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
          <b-field label="Phone number (optional)">
            <b-input
              type="tel"
              v-model="phoneNumber"
              placeholder="Enter your phone number"
              icon="phone"
              size="is-medium"
              maxlength="30">
            </b-input>
          </b-field>

          <b-field label="Birthday (optional)">
            <b-datepicker
              v-model="jsBirthday"
              :max-date="maxDate"
              :min-date="minDate"
              :show-week-numbers="true"
              :focused-date="focusedDate"
              placeholder="Click to select birthday"
              icon="cake-variant"
              size="is-medium"
              trap-focus>
            </b-datepicker>
          </b-field>
          <br/>

          <div class="field has-addons is-marginless">
            <h1 class="title is-6 is-marginless">Password</h1>
            <p class="control">
              <infotooltip message="Make sure that your password has: 10 or more characters, at least one lower case letter, at least one upper case letter, at least one number"/>
            </p>
          </div>
          <b-field>
            <b-input
              type="password"
              v-model="password"
              password-reveal
              maxlength="70"
              placeholder="Enter your password"
              icon="textbox-password"
              size="is-medium"
              pattern="^([a-z]*)([a-z]*).{10,}$"
              validation-message="password is invalid. passwords must include: one number, one lowercase letter, one uppercase letter, and be eight or more characters."
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
              pattern="^([a-z]*)([A-Z]*).{10,}$"
              validation-message="Password is invalid. Passwords must include: one number, one lowercase letter, one uppercase letter, and be eight or more characters."
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
    const today = new Date()
    const maxDate = new Date(today.getFullYear() - 15, today.getMonth(), today.getDay())
    const minDate = new Date(today.getFullYear() - 100, today.getMonth(), today.getDate())

    return {
      minDate: minDate,
      maxDate: maxDate,
      focusedDate: maxDate,
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
    headerDisplay: () => import('@/frontend/components/common/header-display'),
    infotooltip: () => import('@/frontend/components/common/info-tooltip.vue')
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
          common.DisplayFailureToast(resp.data.metadata.response)
          return
        }
        localStorage.setItem('authToken', resp.data.data)
        common.DisplaySuccessToast(resp.data.metadata.response)
        setTimeout(() => {
          loadingComponent.close()
          window.location.href = '/'
        }, 2 * 1000)
      }).catch(err => {
        loadingComponent.close()
        common.DisplayFailureToast(err.response.data.metadata.response)
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
