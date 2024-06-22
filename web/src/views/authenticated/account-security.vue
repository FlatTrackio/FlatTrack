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
    <div class="container">
      <section class="section">
        <nav
          class="breadcrumb is-medium has-arrow-separator"
          aria-label="breadcrumbs"
        >
          <ul>
            <li>
              <router-link :to="{ name: 'Account' }">My account</router-link>
            </li>
            <li class="is-active">
              <router-link :to="{ name: 'Account Security' }"
                >Security</router-link
              >
            </li>
          </ul>
          <b-button
            @click="CopyHrefToClipboard()"
            icon-left="content-copy"
            size="is-small"
          ></b-button>
        </nav>
        <h1 class="title is-1">Security</h1>
        <p class="subtitle is-3">Manage your account's security</p>
        <br />
        <div class="field has-addons">
          <h1 class="title is-3">Password</h1>
          <p class="control">
            <infotooltip
              message="Make sure that your password has: 10 or more characters, at least one lower case letter, at least one upper case letter, at least one number"
            />
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
            pattern="^([a-zA-Z]*).{10,}$"
            validation-message="password is invalid. Make sure that your password has: 10 or more characters, at least one lower case letter, at least one upper case letter, at least one number"
            @keyup.enter.native="PatchProfile"
            maxlength="70"
          >
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
            pattern="^([a-zA-Z]*).{10,}$"
            validation-message="password is invalid. Make sure that your password has: 10 or more characters, at least one lower case letter, at least one upper case letter, at least one number"
            @keyup.enter.native="PatchProfile"
            icon="textbox-password"
          >
          </b-input>
        </b-field>
        <b-button
          type="is-success"
          size="is-medium"
          icon-left="delta"
          native-type="submit"
          expanded
          @click="PatchProfile"
        >
          Update password
        </b-button>

        <br />
        <br />
        <h1 class="title is-3">Two-factor authentication</h1>
        <div class="field">
          <b-checkbox
            :type="otpIsEnabled === true ? 'is-success' : 'is-warning'"
            size="is-medium"
            :indeterminate="otpEnable"
            disabled
            v-model="otpEnable"
          >
            OTP
          </b-checkbox>
        </div>
        <div v-if="otpEnable">
          <qrcode-vue :value="''" :size="200" level="H"></qrcode-vue>
          <br />
          <b-field label="Confirm your OTP code">
            <b-input
              type="number"
              placeholder="Enter your OTP code from your phone"
              maxlength="6"
              size="is-medium"
              icon="qrcode"
            >
            </b-input>
          </b-field>
          <b-button
            type="is-success"
            size="is-medium"
            rounded
            expanded
            native-type="submit"
          >
            Enable
          </b-button>
        </div>

        <br />
        <br />
        <h1 class="title is-3">Sign out of all devices</h1>
        <div class="notification is-warning mb-4">
          <p class="subtitle is-6">
            <strong>Please note:</strong> revoking access is not unable and will
            require signing in again (including from this device).
          </p>
        </div>
        <b-button
          type="is-danger"
          size="is-medium"
          icon-left="close"
          expanded
          @click="ResetAuth"
        >
          Revoke access for all devices
        </b-button>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/common/common'
import profile from '@/requests/authenticated/profile'
import { DialogProgrammatic as Dialog } from 'buefy'

export default {
  name: 'account-security',
  data () {
    const today = new Date()
    const maxDate = new Date(
      today.getFullYear() - 15,
      today.getMonth(),
      today.getDate()
    )
    const minDate = new Date(
      today.getFullYear() - 100,
      today.getMonth(),
      today.getDate()
    )

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
      creationTimestamp: ''
    }
  },
  methods: {
    CopyHrefToClipboard () {
      common.CopyHrefToClipboard()
    },
    GetProfile () {
      profile.GetProfile().then((resp) => {
        this.names = resp.data.spec.names
        this.birthday = resp.data.spec.birthday || null
        this.jsBirthday =
          this.birthday !== null ? new Date(this.birthday * 1000) : null
        this.focusedDate = this.jsBirthday
        this.phoneNumber = resp.data.spec.phoneNumber
        this.groups = resp.data.spec.groups
        this.email = resp.data.spec.email
        this.creationTimestamp = resp.data.spec.creationTimestamp
      })
    },
    PatchProfile () {
      if (this.password !== this.passwordConfirm) {
        common.DisplayFailureToast(
          'Unable to use password as they either do not match'
        )
        return
      }
      profile
        .PatchProfile(
          this.names,
          this.email,
          this.phoneNumber,
          this.birthday,
          this.password
        )
        .then((resp) => {
          if (resp.data.spec.id === '') {
            common.DisplayFailureToast('Failed to update profile')
            return
          }
          common.DisplaySuccessToast('Successfully updated your profile')
        })
        .catch((err) => {
          common.DisplayFailureToast(
            'Failed to update profile' +
              '<br/>' +
              err.response.data.metadata.response
          )
        })
    },
    ResetAuth () {
      Dialog.confirm({
        title: 'Revoke access for all devices',
        message:
          'Are you sure that you wish to sign out of all devices?' +
          '<br/>' +
          'This action cannot be undone.',
        confirmText: 'Sign out',
        type: 'is-danger',
        hasIcon: true,
        onConfirm: () => {
          profile
            .PostAuthReset()
            .then((resp) => {
              common.DisplaySuccessToast(
                'Successfully signed out of all devices'
              )
              window.location.href = '/login'
            })
            .catch((err) => {
              common.DisplayFailureToast(
                'Failed to sign out of all devices' +
                  '<br/>' +
                  err.response.data.metadata.response
              )
            })
        }
      })
    },
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    }
  },
  computed: {
    birthday () {
      return new Date(this.jsBirthday || 0).getTime() / 1000 || 0
    }
  },
  components: {
    QrcodeVue: () => import('qrcode.vue'),
    infotooltip: () => import('@/components/common/info-tooltip.vue')
  }
}
</script>
