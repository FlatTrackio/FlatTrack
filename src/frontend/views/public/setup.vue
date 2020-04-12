<template>
<div>
  <headerDisplay/>
  <div class="container">
    <section class="section form-width">
      <h1 class="title is-1">Welcome to FlatTrack!</h1>

      <h1 class="title is-1">Set up</h1>
      <p class="subtitle is-5">
        Let's get started and set up your instance.
      </p>
      <br/>
      <div class="">
        <b-icon
          icon="cogs"
          size="is-medium">
        </b-icon>
        <h3 class="title is-4">System</h3>
        <b-field label="Language">
          <b-select
            placeholder="English"
            v-model="language"
            autofocus
            required
            icon="web"
            size="is-medium"
            expanded>
            <option value="English">English</option>
          </b-select>
        </b-field>
        <b-field label="Timezone">
          <b-select
            placeholder="Pacific/Auckland"
            v-model="timezone"
            autofocus
            required
            icon="map-clock"
            size="is-medium"
            expanded>
            <option value="Pacific/Auckland">Pacific/Auckland</option>
          </b-select>
        </b-field>
        <br/>
        <b-icon
          icon="home"
          size="is-medium">
        </b-icon>

        <h3 class="title is-4">Your flat</h3>
        <b-field label="Flat name">
          <b-input type="text"
                   v-model="flatName"
                   maxlength="20"
                   placeholder="Enter what your flat is referred to as"
                   icon="textbox"
                   size="is-medium"
                   required>
          </b-input>
        </b-field>
        <b-icon
          icon="account"
          size="is-medium">
        </b-icon>

        <h3 class="title is-4">Your profile</h3>
        <p class="subtitle is-6">Note: your account profile will be set up as Administrator</p>
        <b-field label="Name(s)">
          <b-input type="text"
                   v-model="names"
                   maxlength="70"
                   placeholder="Enter your name(s)"
                   icon="textbox"
                   size="is-medium"
                   required>
          </b-input>
        </b-field>
        <b-field label="Email">
          <b-input type="email"
                   v-model="email"
                   maxlength="70"
                   placeholder="Enter your email address"
                   icon="email"
                   size="is-medium"
                   required>
          </b-input>
        </b-field>
        <b-field label="Phone number*">
          <b-input type="tel"
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
          <b-input type="password"
                   v-model="password"
                   password-reveal
                   maxlength="70"
                   placeholder="Enter a password"
                   icon="textbox-password"
                   size="is-medium"
                   required>
          </b-input>
        </b-field>
        <b-field label="Confirm password">
          <b-input type="password"
                   v-model="passwordConfirm"
                   password-reveal
                   placeholder="Confirm your password"
                   icon="textbox-password"
                   @keyup.enter.native="Register({ language, timezone, flatName, user: { names, email, password } })"
                   size="is-medium"
                   maxlength="70">
          </b-input>
        </b-field>
        <b-button type="is-success" size="is-medium" rounded native-type="submit" @click="Register({ language, timezone, flatName, user: { names, email, password, passwordConfirm, jsBirthday, phoneNumber } })">Setup</b-button>
      </div>
      <br/>
      <p>* optional</p>
      <br/>
      <p>FlatTrack version {{ flatTrackVersion || '0.0.0' }}</p>
    </section>
  </div>
</div>
</template>

<script>
import headerDisplay from '@/frontend/components/header-display.vue'
import registration from '@/frontend/requests/public/registration'
import apiroot from '@/frontend/requests/public/apiroot'
import common from '@/frontend/common/common'
import { LoadingProgrammatic as Loading } from 'buefy'
import moment from 'moment'

export default {
  name: 'setup',
  data () {
    const today = new Date()
    const maxDate = new Date(today.getFullYear() - 15, today.getMonth(), today.getDay())

    return {
      flatTrackVersion: '0.0.0',
      language: 'English',
      timezone: 'Pacific/Auckland',
      jsBirthday: null,
      passwordConfirm: '',
      flatName: null,
      names: null,
      email: null,
      phoneNumber: null,
      birthday: null,
      password: null
    }
  },
  components: {
    headerDisplay
  },
  methods: {
    Register: (form) => {
      if (form.user.password !== form.user.passwordConfirm) {
        common.DisplayFailureToast('Error passwords do not match')
        return
      }
      form.user.birthday = Number(moment(form.user.jsBirthday).format('X')) || 0

      const loadingComponent = Loading.open({
        container: null
      })
      setTimeout(() => loadingComponent.close(), 20 * 1000)
      registration.PostAdminRegister(form).then(resp => {
        if (resp.data.data !== '' || typeof resp.data.data !== 'undefined') {
          localStorage.setItem('authToken', resp.data.data)
        } else {
          Error('failed to find authToken')
        }
        common.DisplaySuccessToast('Welcome to FlatTrack!')
        setTimeout(() => {
          loadingComponent.close()
          window.location.href = '/'
        }, 3 * 1000)
      }).catch(err => {
        loadingComponent.close()
        common.DisplayFailureToast(err.response.data.metadata.response)
      })
    },
    GetAPIversion () {
      apiroot.GetAPIroot().then(resp => {
        this.flatTrackVersion = resp.data.metadata.version
      })
    }
  },
  async beforeMount () {
    this.GetAPIversion()
  }
}
</script>

<style scoped>
.form-width {
    width: 380px;
    margin: auto;
}
</style>
