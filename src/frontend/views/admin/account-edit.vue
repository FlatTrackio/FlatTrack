<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><router-link to="/admin/accounts">Accounts</router-link></li>
              <li class="is-active"><router-link to="/admin/accounts/new">Edit account</router-link></li>
            </ul>
        </nav>
        <h1 class="title is-1">Edit account</h1>
        <p class="subtitle is-3">Edit an existing user account</p>

        <div v-if="registered !== true">
          <div class="notification is-warning">
            <p class="subtitle is-6">This account doesn't appear to be registered</p>
            <b-button @click="showRegistrationCompletionDetails = !showRegistrationCompletionDetails">{{ showRegistrationCompletionDetails === false ? 'Show' : 'Hide' }} registration details</b-button>
            <div v-if="showRegistrationCompletionDetails === true">
              <br/>
              <div class="notification">
                <div class="content">
                  <qrcode-vue :value="windowOrigin + '/useraccountconfirm/' + userAccountConfirmId + '?secret=' + userAccountConfirmSecret" :size="200" level="H"></qrcode-vue>
                  <br/>
                  <p>
                    Have your flatmate scan the QR code above, or <a type="is-text" @click="CopyRegistrationLink">click here</a> to copy the registration link for you to send to your flatmate
                  </p>
                </div>
              </div>
            </div>
          </div>
          <br/>
        </div>

        <b-field label="Name(s)">
          <b-input type="text"
                   v-model="names"
                   maxlength="60"
                   placeholder="Enter your flatmate's name"
                   icon="textbox"
                   size="is-medium"
                   required>
          </b-input>
        </b-field>

        <b-field label="Email">
          <b-input type="email"
                   v-model="email"
                   maxlength="70"
                   placeholder="Enter your flatmate's email"
                   icon="email"
                   size="is-medium"
                   required>
          </b-input>
        </b-field>

        <section>
          <b-field label="Groups">
            <b-taginput
              v-model="groupsFull"
              :data="availableGroups"
              field="name"
              autocomplete
              open-on-focus
              ellipsis
              icon="account-group"
              placeholder="Select groups"
              size="is-medium"
              @typing="GetFilteredGroups"
              >
          </b-field>
        </section>
        <br/>

        <b-field label="Phone number">
          <b-input type="tel"
                   v-model="phoneNumber"
                   placeholder="Enter your flatmate's phone number"
                   icon="phone"
                   size="is-medium"
                   maxlength="30">
          </b-input>
        </b-field>

        <b-field label="Birthday">
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
                   placeholder="Enter a password for your flatmate"
                   icon="textbox-password"
                   size="is-medium"
                   required>
          </b-input>
        </b-field>

        <b-field label="Confirm password">
          <b-input type="password"
                   v-model="passwordConfirm"
                   password-reveal
                   maxlength="70"
                   placeholder="Confirm a password for your flatmate"
                   icon="textbox-password"
                   size="is-medium"
                   required>
          </b-input>
        </b-field>
        <b-button type="is-success" size="is-medium" rounded native-type="submit" @click="PatchUserAccount(names, email, phoneNumber, birthday, password, passwordConfirm, jsBirthday, groupsFull)">Update user account</b-button>
        <b-button type="is-danger" size="is-medium" rounded native-type="submit" @click="DeleteUserAccount(id)">Delete user account</b-button>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/frontend/common/common'
import groups from '@/frontend/requests/authenticated/groups'
import flatmates from '@/frontend/requests/authenticated/flatmates'
import adminFlatmates from '@/frontend/requests/admin/flatmates'
import moment from 'moment'
import { DialogProgrammatic as Dialog } from 'buefy'

export default {
  name: 'new account',
  data () {
    var today = new Date()
    var maxDate = new Date(today.getFullYear() - 15, today.getMonth(), today.getYear())
    var windowOrigin = window.location.origin

    return {
      windowOrigin: windowOrigin,
      maxDate: maxDate,
      showRegistrationCompletionDetails: false,
      userAccountConfirmId: null,
      userAccountConfirmSecret: null,
      id: this.$route.params.id,
      names: null,
      email: null,
      phoneNumber: null,
      birthday: 0,
      groups: [],
      registered: null,
      password: null,
      passwordConfirm: null,
      availableGroups: [],
      jsBirthday: null,
      groupsFull: []
    }
  },
  components: {
    QrcodeVue: () => import('qrcode.vue')
  },
  methods: {
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    },
    GetAvailableGroups () {
      groups.GetGroups().then(resp => {
        this.availableGroups = resp.data.list
        resp.data.list.map(group => {
          if (group.defaultGroup === true) {
            this.groupsFull.push(group)
          }
        })
      }).catch(err => {
        console.log(err)
        common.DisplayFailureToast('Failed to list groups')
      })
    },
    GetFilteredGroups (text) {
      this.groups = this.availableGroups.filter((group) => {
        return group.name
          .toString()
          .toLowerCase()
          .indexOf(text.toLowerCase()) >= 0
      })
    },
    GetUserAccount () {
      flatmates.GetFlatmate(this.id).then(resp => {
        var user = resp.data.spec
        this.names = user.names
        this.email = user.email
        this.phoneNumber = user.phoneNumber
        this.birthday = user.birthday
        this.registered = user.registered
        this.groups = user.groups
        this.groupsFull = []
        this.availableGroups.map(group => {
          if (this.groups.includes(group.name)) {
            this.groupsFull.push(group)
          }
        })
        this.jsBirthday = typeof this.birthday !== 'undefined' ? new Date(this.birthday * 1000) : null
      }).catch(err => {
        console.log({ err })
        common.DisplayFailureToast('Failed to fetch user account' + '<br/>' + (err.response.data.metadata.response || err))
      })
    },
    PatchUserAccount (names, email, phoneNumber, birthday, password, passwordConfirm, jsBirthday, groupsFull) {
      if (password !== passwordConfirm && password !== null && typeof password !== 'undefined') {
        common.DisplayFailureToast('Passwords do not match')
        return
      }
      birthday = Number(moment(jsBirthday).format('X'))
      var groups = []
      groupsFull.map(group => {
        if (group === '' || group.name === '') {
          return
        }
        groups.push(group.name)
      })
      adminFlatmates.PatchFlatmate(this.id, names, email, phoneNumber, birthday, groups, password).then(resp => {
        common.DisplaySuccessToast('Updated user account')
        setTimeout(() => {
          this.$router.push({ name: 'Admin accounts' })
        }, 1 * 1000)
      }).catch(err => {
        common.DisplayFailureToast('Failed to create user account' + '<br/>' + (err.response.data.metadata.response || err))
      })
    },
    DeleteUserAccount (id) {
      Dialog.confirm({
        title: 'Delete user account',
        message: 'Are you sure that you wish to remove this account?' + '<br/>' + 'This action cannot be undone.',
        confirmText: 'Delete account',
        type: 'is-danger',
        hasIcon: true,
        onConfirm: () => {
          adminFlatmates.DeleteFlatmate(id).then(resp => {
            common.DisplaySuccessToast('Deleted user account')
            setTimeout(() => {
              this.$router.push({ name: 'Admin accounts' })
            }, 1 * 1000)
          }).catch(err => {
            common.DisplayFailureToast('Failed to delete user account' + '<br/>' + (err.response.data.metadata.response || err))
          })
        }
      })
    },
    CopyRegistrationLink () {
      var registrationLink = `${window.location.origin}/useraccountconfirm/${this.userAccountConfirmId}?secret=${this.userAccountConfirmSecret}`
      window.prompt('Copy the following link', registrationLink)
    }
  },
  async beforeMount () {
    this.GetAvailableGroups()
    this.GetUserAccount()
    if (this.registered !== true) {
      adminFlatmates.GetUserAccountConfirms(this.id).then(resp => {
        var confirmsList = resp.data.list
        for (var confirmItem in confirmsList) {
          if (confirmsList[confirmItem].userId !== this.id) {
            continue
          }
          this.userAccountConfirmId = confirmsList[confirmItem].id
          this.userAccountConfirmSecret = confirmsList[confirmItem].secret
        }
      })
    }
  }
}
</script>
