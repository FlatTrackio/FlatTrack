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
            @keyup.enter.native="Register"
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
            @keyup.enter.native="Register"
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
          <b-input
            type="text"
            v-model="flatName"
            maxlength="20"
            placeholder="Enter your flat's name"
            icon="textbox"
            size="is-medium"
            @keyup.enter.native="Register"
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
          <b-input
            type="text"
            v-model="names"
            maxlength="70"
            placeholder="Enter your name(s)"
            icon="textbox"
            size="is-medium"
            @keyup.enter.native="Register"
            required>
          </b-input>
        </b-field>
        <b-field label="Email">
          <b-input
            type="email"
            v-model="email"
            maxlength="70"
            placeholder="Enter your email address"
            icon="email"
            size="is-medium"
            @keyup.enter.native="Register"
            required>
          </b-input>
        </b-field>
        <b-field label="Phone number (optional)">
          <b-input
            type="tel"
            v-model="phoneNumber"
            placeholder="Enter your phone number"
            icon="phone"
            size="is-medium"
            @keyup.enter.native="Register"
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
            @keyup.enter.native="Register"
            trap-focus>
          </b-datepicker>
        </b-field>
        <br/>
        <div class="field has-addons">
          <label class="label">Password</label>
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
            placeholder="Enter a password"
            icon="textbox-password"
            size="is-medium"
            pattern="^([a-z]*)([A-Z]*).{10,}$"
            validation-message="Password is invalid. Passwords must include: one number, one lowercase letter, one uppercase letter, and be eight or more characters."
            @keyup.enter.native="Register"
            required>
          </b-input>
        </b-field>
        <b-field label="Confirm password">
          <b-input
            type="password"
            v-model="passwordConfirm"
            password-reveal
            placeholder="Confirm your password"
            icon="textbox-password"
            @keyup.enter.native="Register"
            size="is-medium"
            maxlength="70"
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
          expanded
          @click="Register">
          Setup
        </b-button>
      </div>
    </section>
  </div>
</div>
</template>

<script>
import registration from '@/frontend/requests/public/registration'
import apiroot from '@/frontend/requests/public/apiroot'
import common from '@/frontend/common/common'
import { LoadingProgrammatic as Loading } from 'buefy'

export default {
  name: 'setup',
  data () {
    const today = new Date()
    const maxDate = new Date(today.getFullYear() - 15, today.getMonth(), today.getDay())
    const minDate = new Date(today.getFullYear() - 100, today.getMonth(), today.getDay())
    const focusedDate = new Date(today.getFullYear() - 15, today.getMonth() - 5, today.getDay())

    return {
      language: 'English',
      timezone: 'Pacific/Auckland',
      maxDate: maxDate,
      minDate: minDate,
      focusedDate: focusedDate,
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
    headerDisplay: () => import('@/frontend/components/common/header-display.vue'),
    infotooltip: () => import('@/frontend/components/common/info-tooltip.vue')
  },
  methods: {
    Register () {
      if (this.password !== this.passwordConfirm) {
        common.DisplayFailureToast('Error passwords do not match')
        return
      }

      this.birthday = new Date(this.jsBirthday || 0).getTime() / 1000 || 0

      const loadingComponent = Loading.open({
        container: null
      })
      setTimeout(() => loadingComponent.close(), 20 * 1000)
      var form = {
        language: this.language,
        timezone: this.timezone,
        flatName: this.flatName,
        user: {
          names: this.names,
          email: this.email,
          password: this.password,
          passwordConfirm: this.passwordConfirm,
          jsBirthday: this.jsBirthday,
          phoneNumber: this.phoneNumber
        }
      }
      registration.PostAdminRegister(form).then(resp => {
        if (resp.data.data !== '' || typeof resp.data.data !== 'undefined') {
          localStorage.setItem('authToken', resp.data.data)
        } else {
          common.DisplayFailureToast('Failed to find login token after registration')
          return
        }
        common.DisplaySuccessToast('Welcome to FlatTrack!')
        setTimeout(() => {
          loadingComponent.close()
          this.$router.push({ name: 'Home' })
        }, 3 * 1000)
      }).catch(err => {
        loadingComponent.close()
        common.DisplayFailureToast(err.response.data.metadata.response || err)
      })
    }
  }
}
</script>

<style scoped>
</style>
