<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><router-link to="/admin">Admin</router-link></li>
              <li><router-link to="/admin/accounts">Accounts</router-link></li>
              <li class="is-active"><router-link to="/admin/accounts/new">New account</router-link></li>
            </ul>
        </nav>
        <h1 class="title is-1">Edit account</h1>
        <p class="subtitle is-3">Edit an existing user account</p>

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
        <b-button type="is-success" size="is-medium" rounded disabled native-type="submit" @click="UpdateUserAccount(names, email, phoneNumber, birthday, groups, password, passwordConfirm, jsBirthday, groupsFull)">Update user account</b-button>
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
    return {
      id: this.$route.params.id,
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
    GetUserAccount () {
      flatmates.GetFlatmate(this.id).then(resp => {
        var user = resp.data.spec
        this.names = user.names
        this.email = user.email
        this.phoneNumber = user.phoneNumber
        this.birthday = user.birthday
        this.groups = user.groups
      }).catch(err => {
        console.log({ err })
        common.DisplayFailureToast('Failed to fetch user account' + err.response.data.metadata.response)
      })
    },
    UpdateUserAccount (names, email, phoneNumber, birthday, groups, password, passwordConfirm, jsBirthday, groupsFull) {
      if (password !== passwordConfirm && password !== '') {
        common.DisplayFailureToast('Passwords do not match')
      }
      birthday = moment(jsBirthday).format('X')
      groupsFull.map(group => {
        groups = [...groups, group.name]
      })
      adminFlatmates.PostFlatmate({ names, email, phoneNumber, birthday, groups, password }).then(resp => {
        common.DisplaySuccessToast('Created user account')
        setTimeout(() => {
          this.$router.push({ name: 'Admin accounts' })
        }, 1 * 1000)
      }).catch(err => {
        common.DisplayFailureToast('Failed to create user account' + `<br/>${err.response.data.metadata.response}`)
      })
    },
    DeleteUserAccount (id) {
      Dialog.confirm({
        title: 'Delete user account',
        message: `Are you sure that you wish to remove this account?` + '<br/>' + `This action cannot be undone.`,
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
            common.DisplayFailureToast('Failed to delete user account' + `<br/>${err.response.data.metadata.response}`)
          })
        }
      })
    }
  },
  async created () {
    this.GetAvailableGroups()
    this.GetUserAccount()
  }
}
</script>
