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
        <div class="field has-addons">
          <h1 class="title is-3">Password</h1>
          <p class="control">
            <infotooltip message="Make sure that your password has: 10 or more characters, at least one lower case letter, at least one upper case letter, at least one number"/>
          </p>
        </div>
        <b-field label="Password">
          <b-input
            type="password"
            v-model="password"
            password-reveal
            placeholder="Enter to update your password"
            icon="textbox-password"
            size="is-medium"
            pattern="^([a-z]*)([A-Z]*).{10,}$"
            validation-message="Password is invalid. Passwords must include: one number, one lowercase letter, one uppercase letter, and be eight or more characters."
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
            pattern="^([a-z]*)([A-Z]*).{10,}$"
            validation-message="Password is invalid. Passwords must include: one number, one lowercase letter, one uppercase letter, and be eight or more characters."
            icon="textbox-password">
          </b-input>
        </b-field>
        <b-button
          type="is-success"
          size="is-medium"
          icon-left="delta"
          native-type="submit"
          @click="PatchProfile(names, email, phoneNumber, password, passwordConfirm, jsBirthday, resourceVersion)">
          Update password
        </b-button>

        <br/>
        <br/>
        <h1 class="title is-3">Two-factor authentication</h1>
        <div class="field">
          <b-checkbox
            :type="otpIsEnabled === true ? 'is-success' : 'is-warning'"
            size="is-medium"
            :indeterminate="otpEnable"
            disabled
            v-model="otpEnable">
            OTP
          </b-checkbox>
        </div>
        <div v-if="otpEnable">
          <qrcode-vue :value="''" :size="200" level="H"></qrcode-vue>
          <br/>
          <b-field label="Confirm your OTP code">
            <b-input
              type="number"
              placeholder="Enter your OTP code from your phone"
              maxlength="6"
              size="is-medium"
              icon="qrcode">
            </b-input>
          </b-field>
          <b-button
            type="is-success"
            size="is-medium"
            rounded
            native-type="submit">
            Enable
          </b-button>
        </div>

        <br/>
        <br/>
        <h1 class="title is-3">Sign out of all devices</h1>
        <b-button
          type="is-danger"
          size="is-medium"
          icon-left="close"
          @click="ResetAuth">
          Revoke access for all devices
        </b-button>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/frontend/common/common'
import profile from '@/frontend/requests/authenticated/profile'
import { DialogProgrammatic as Dialog } from 'buefy'

export default {
  name: 'account security',
  data () {
    const today = new Date()
    const maxDate = new Date(today.getFullYear() - 15, today.getMonth(), today.getDate())
    const minDate = new Date(today.getFullYear() - 100, today.getMonth(), today.getDate())

    return {
      maxDate: maxDate,
      minDate: minDate,
      focusedDate: maxDate,
      passwordConfirm: '',
      otpEnable: false,
      otpIsEnabled: false,
      jsBirthday: null,
      names: '',
      email: '',
      phoneNumber: '',
      groups: [],
      password: '',
      birthday: 0,
      creationTimestamp: '',
      resourceVersion: ''
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
        this.resourceVersion = resp.data.spec.resourceVersion
      })
    },
    PatchProfile (names, email, phoneNumber, password, passwordConfirm, jsBirthday, resourceVersion) {
      if (password !== passwordConfirm) {
        common.DisplayFailureToast('Unable to use password as they either do not match')
        return
      }
      var birthday = new Date(jsBirthday || 0).getTime() / 1000 || 0
      profile.PatchProfile(names, email, phoneNumber, birthday, password, resourceVersion).then(resp => {
        if (resp.data.spec.id === '') {
          common.DisplayFailureToast('Failed to update profile')
          return
        }
        this.resourceVersion = resp.data.spec.resourceVersion
        common.DisplaySuccessToast('Successfully updated your profile')
      }).catch(err => {
        common.DisplayFailureToast('Failed to update profile' + '<br/>' + err.response.data.metadata.response)
      })
    },
    ResetAuth () {
      Dialog.confirm({
        title: 'Revoke access for all devices',
        message: 'Are you sure that you wish to sign out of all devices?' + '<br/>' + 'This action cannot be undone.',
        confirmText: 'Sign out',
        type: 'is-danger',
        hasIcon: true,
        onConfirm: () => {
          profile.PostAuthReset().then(resp => {
            common.DisplaySuccessToast('Successfully signed out of all devices')
            window.location.href = '/login'
          }).catch(err => {
            common.DisplayFailureToast('Failed to sign out of all devices' + '<br/>' + err.response.data.metadata.response)
          })
        }
      })
    },
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    }
  },
  components: {
    QrcodeVue: () => import('qrcode.vue'),
    infotooltip: () => import('@/frontend/components/common/info-tooltip.vue')
  }
}
</script>
