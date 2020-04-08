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
                   required>
          </b-input>
        </b-field>

        <b-field label="Email">
          <b-input type="email"
                   v-model="email"
                   maxlength="70"
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
              icon="label"
              placeholder="Select groups"
              @typing="GetFilteredGroups"
              >
          </b-field>
        </section>
        <br/>

        <b-field label="Phone number">
          <b-input type="tel"
                   v-model="phoneNumber"
                   maxlength="30"
                   >
          </b-input>
        </b-field>

        <b-field label="Birthday">
          <b-datepicker
            v-model="jsBirthday"
            show-week-number
            inline
            placeholder="Click to select birthday"
            icon="calendar-today"
            trap-focus>
            </b-datepicker>
        </b-field>
        <br/>

        <b-field label="Password">
          <b-input type="password"
                   v-model="password"
                   password-reveal
                   maxlength="70"
                   required>
          </b-input>
        </b-field>

        <b-field label="Confirm password">
          <b-input type="password"
                   v-model="passwordConfirm"
                   password-reveal
                   maxlength="70"
                   required>
          </b-input>
        </b-field>
        <b-button type="is-success" size="is-medium" rounded native-type="submit" @click="PostNewUser(names, email, phoneNumber, birthday, groups, password, passwordConfirm, jsBirthday, groupsFull)">Create user account</b-button>
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
    return {
      names: '',
      email: '',
      phoneNumber: '',
      birthday: 0,
      groups: [],
      password: '',
      passwordConfirm: '',
      availableGroups: [],
      jsBirthday: new Date(),
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
    PostNewUser (names, email, phoneNumber, birthday, groups, password, passwordConfirm, jsBirthday, groupsFull) {
      if (password !== passwordConfirm && password !== '') {
        common.DisplayFailureToast('Passwords do not match')
      }
      birthday = moment(jsBirthday).format('X')
      groupsFull.map(group => {
        groups = [...groups, group.name]
      })
      console.log({ names, email, phoneNumber, birthday, groups, password, passwordConfirm, jsBirthday, groupsFull })
      console.log({ adminFlatmates })
      adminFlatmates.PostFlatmate({ names, email, phoneNumber, birthday, groups, password }).then(resp => {
        common.DisplaySuccessToast('Created user account')
        setTimeout(() => {
          this.$router.push({ name: 'Admin accounts' })
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
