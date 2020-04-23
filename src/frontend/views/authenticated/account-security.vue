<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><router-link to="/account">My account</router-link></li>
              <li class="is-active"><router-link to="/account/security">Security</router-link></li>
            </ul>
        </nav>
        <h1 class="title is-1">Security</h1>
        <p class="subtitle is-3">Manage your account's security</p>
        <br/>
        <h1 class="title is-3">Password<h1/>
        <b-field label="Password">
          <b-input
            type="password"
            v-model="password"
            password-reveal
            placeholder="Enter to update your password"
            icon="textbox-password"
            size="is-medium"
            maxlength="70">
          </b-input>
        </b-field>
        <b-field label="Confirm password">
          <b-input
            type="password"
            v-model="passwordConfirm"
            password-reveal
            placeholder="Confirm to update your password"
            maxlength="70"
            size="is-medium"
            icon="textbox-password">
          </b-input>
        </b-field>
        <b-button
          type="is-success"
          size="is-medium"
          rounded
          native-type="submit"
          @click="PatchProfile(names, email, phoneNumber, password, passwordConfirm, jsBirthday)">
          Update profile
        </b-button>

        <br/>
        <br/>
        <h1 class="title is-3">Two-factor authentication</h1>
        <div class="field">
          <b-checkbox size="is-medium">OTP</b-checkbox>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/frontend/common/common'
import profile from '@/frontend/requests/authenticated/profile'

export default {
  name: 'profile',
  data () {
    const today = new Date()
    const maxDate = new Date(today.getFullYear() - 15, today.getMonth(), today.getDate())
    const minDate = new Date(today.getFullYear() - 100, today.getMonth(), today.getDate())

    return {
      maxDate: maxDate,
      minDate: minDate,
      focusedDate: maxDate,
      passwordConfirm: '',
      jsBirthday: null,
      names: '',
      email: '',
      phoneNumber: '',
      groups: [],
      password: '',
      birthday: 0,
      creationTimestamp: ''
    }
  },
  methods: {
    GetProfile () {
      profile.GetProfile().then(resp => {
        this.names = resp.data.spec.names
        this.birthday = resp.data.spec.birthday || null
        this.jsBirthday = this.birthday !== null ? new Date(this.birthday * 1000) : null
        this.focusedDate = this.jsBirthday
        this.phoneNumber = resp.data.spec.phoneNumber
        this.groups = resp.data.spec.groups
        this.email = resp.data.spec.email
        this.creationTimestamp = resp.data.spec.creationTimestamp
      })
    },
    PatchProfile (names, email, phoneNumber, password, passwordConfirm, jsBirthday) {
      if (password !== passwordConfirm) {
        common.DisplayFailureToast('Unable to use password as they either do not match')
        return
      }
      var birthday = new Date(jsBirthday || 0).getTime() / 1000 || 0
      profile.PatchProfile(names, email, phoneNumber, birthday, password).then(resp => {
        if (resp.data.spec.id === '') {
          common.DisplayFailureToast('Failed to update profile')
          return
        }
        common.DisplaySuccessToast('Successfully updated your profile')
      }).catch(err => {
        common.DisplayFailureToast('Failed to update profile' + '<br/>' + err.response.data.metadata.response)
      })
    },
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    }
  },
  async beforeMount () {
    this.GetProfile()
  }
}
</script>
