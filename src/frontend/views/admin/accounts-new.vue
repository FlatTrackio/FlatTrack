<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><router-link to="/admin/accounts">Accounts</router-link></li>
              <li class="is-active"><router-link to="/admin/accounts/new">New account</router-link></li>
            </ul>
        </nav>
        <h1 class="title is-1">New account</h1>
        <p class="subtitle is-3">Add a new flatmate</p>

        <b-field label="Names">
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
              @typing="GetFilteredGroups">
          </b-field>
        </section>
        <br/>

        <div class="field">
          <b-checkbox v-model="setOnlyRequiredFields">Allow your flatmate to set their password, phone number, and birthday</b-checkbox>
        </div>
        <div v-if="!setOnlyRequiredFields">
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
        </div>
        <div v-else>
          <p class="subtitle is-6"><b>Please note:</b> email account verification does not work yet, however QR code verfication does. If this in an inconvenience, uncheck the checkbox above to fill all fields out for the new account.</p>
          <br/>
        </div>

        <!-- TODO become invite via email button -->
        <b-button type="is-success" size="is-medium" rounded native-type="submit" @click="PostNewUser(names, email, phoneNumber, birthday, password, passwordConfirm, jsBirthday, groupsFull)">Create user account</b-button>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/frontend/common/common'
import groups from '@/frontend/requests/authenticated/groups'
import adminFlatmates from '@/frontend/requests/admin/flatmates'
import moment from 'moment'

export default {
  name: 'new account',
  data () {
    const today = new Date()
    const maxDate = new Date(today.getFullYear() - 15, today.getMonth(), today.getDay())

    return {
      focusedDate: maxDate,
      maxDate: maxDate,
      setOnlyRequiredFields: true,
      names: null,
      email: null,
      phoneNumber: null,
      birthday: 0,
      password: null,
      passwordConfirm: null,
      availableGroups: [],
      jsBirthday: null,
      groupsFull: []
    }
  },
  methods: {
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    },
    GetAvailableGroups () {
      groups.GetGroups().then(resp => {
        this.availableGroups = resp.data.list
        this.groups = resp.data.list
        resp.data.list.map(group => {
          if (group.defaultGroup === true) {
            this.groupsFull = [...this.groupsFull, group]
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
    PostNewUser (names, email, phoneNumber, birthday, password, passwordConfirm, jsBirthday, groupsFull) {
      if (password !== passwordConfirm && password !== '') {
        common.DisplayFailureToast('Passwords do not match')
        return
      }
      birthday = Number(moment(jsBirthday).format('X')) || 0
      var groups = []
      groupsFull.map(group => {
        if (group === '' || group.name === '') {
          return
        }
        groups.push(group.name)
      })
      adminFlatmates.PostFlatmate({ names, email, phoneNumber, birthday, groups, password }).then(resp => {
        common.DisplaySuccessToast('Created user account')
        setTimeout(() => {
          if (this.setOnlyRequiredFields === true) {
            this.$router.push({ name: 'Admin accounts' })
          } else {
            this.$router.push({ path: '/admin/' })
          }
        }, 1.5 * 1000)
      }).catch(err => {
        common.DisplayFailureToast('Failed to create user account' + `<br/>${err.response.data.metadata.response}`)
      })
    }
  },
  async beforeMount () {
    this.GetAvailableGroups()
  }
}
</script>
