<template>
    <div>
      <headerDisplay admin="true"/>
      <div class="container">
        <section class="section">
          <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
                <li><a href="/#/admin">Admin Home</a></li>
                <li><a href="/#/admin/members">Manage Flatmates</a></li>
                <li class="is-active"><a>{{ returnNamesforID(id, names) || 'Add a flatmate' }}</a></li>
            </ul>
          </nav>
          <h1 class="title">{{ returnNamesforID(id, names) || 'Add a flatmate' }}</h1>
          <div class="card">
            <div class="card-content">
              <div class="content">
                <b-field label="Fullname" v-if="!id">
                    <b-input placeholder="xxxxx xxxxxxxx" v-model="names" maxlength="30" rounded required></b-input>
                </b-field>
                <b-field label="Phone Number (optional)">
                    <b-input placeholder="xx xxx xxxx" v-model="phoneNumber" maxlength="30" rounded></b-input>
                </b-field>
                <b-field label="Email">
                    <b-input placeholder="xxxxx@xxxxx.xxx" type="email" v-model="email" maxlength="30" rounded required></b-input>
                </b-field>
                <b-field label="Allergies (optional)">
                    <b-input placeholder="xx, xxxx, xx" v-model="allergies" maxlength="30" rounded></b-input>
                </b-field>
                <b-field label="Password (optional)">
                    <b-input type="password" placeholder="xxxxxxxxxxxxx" v-model="password" maxlength="30" rounded :required="id" password-reveal :disabled="memberSetPassword"></b-input>
                </b-field>
                <div v-if="!id">
                  <b-checkbox v-model="memberSetPassword">
                    Allow the new member to set the password?
                  </b-checkbox>
                  <br/><br/>
                </div>
                <div class="control">
                  <label class="label">Group</label>
                  <div class="field">
                    <b-radio v-model="group"
                        native-value="flatmember"
                        name="group">
                        Flatmember - standard user, can use anything enabled
                    </b-radio>
                  </div>
                  <div class="field">
                    <b-radio v-model="group"
                        native-value="admin"
                        name="group"
                        type="is-danger">
                        Admin - can access the admin pages
                    </b-radio>
                  </div>
                  <div class="field">
                    <b-radio v-model="group"
                        native-value="approver"
                        name="group"
                        type="is-danger">
                        Approver - can view and approve tasks
                    </b-radio>
                  </div>
                </div>
                <br>
                <div v-if="returnNamesforID(id, names) !== 'undefined'">
                  <b-button type="is-success" @click="updateMember(id, names, email, phoneNumber, allergies, password, group, memberSetPassword)">Update</b-button>
                  <b-button type="is-warning">Disable</b-button>
                  <b-button type="is-danger" @click="deleteMember(id)">Delete</b-button>
                </div>
                <div v-else>
                  <b-button type="is-success" native-type="submit" @click="addNewMember(names, email, phoneNumber, allergies, password, group, memberSetPassword)">Add new flatmate</b-button>
                </div>
                <br>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>
</template>

<script>
import axios from 'axios'
import { Service } from 'axios-middleware'
import { ToastProgrammatic as Toast, DialogProgrammatic as Dialog } from 'buefy'
import headerDisplay from '../common/header-display'

const service = new Service(axios)
service.register({
  onResponse (response) {
    if (response.status === 403) {
      localStorage.removeItem('authToken')
      location.href = '/'
    }
    return response
  }
})
export default {
  name: 'Admin home',
  data () {
    return {
      id: this.$route.query.id,
      names: '',
      email: '',
      phoneNumber: '',
      allergies: '',
      password: '',
      memberSetPassword: true,
      group: 'flatmember',
      member: {
      },
      pageErrors: []
    }
  },
  created () {
    var id = this.$route.query.id
    if (!id) {
      return
    }
    axios.get(`/api/members/${id}`,
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('authToken')}`
        }
      })
      .then(resp => {
        var member = resp.data
        this.names = resp.data.names
        this.email = member.email
        this.phoneNumber = member.phoneNumber
        this.allergies = member.allergies
        this.memberSetPassword = member.memberSetPassword
        this.group = member.group
        this.password = null
        this.member = member
      })
      .catch(err => {
        this.pageErrors = [...this.pageErrors, err]
      })
  },
  methods: {
    addNewMember: (names, email, phoneNumber, allergies, password, group, memberSetPassword) => {
      var member = {
        names,
        email,
        phoneNumber,
        allergies,
        password,
        group,
        memberSetPassword
      }
      axios({
        method: 'post',
        url: `/api/members`,
        data: member,
        headers: {
          Authorization: `Bearer ${localStorage.getItem('authToken')}`
        }}).then(resp => {
        console.log('Add successful', resp)
        Toast.open({
          message: 'Flatmate added successfully',
          position: 'is-bottom',
          type: 'is-success'
        })
        location.href = '#/admin/members'
      })
        .catch(err => {
          console.log('Add failed', err)
          Toast.open({
            message: 'Failed to add flatmate',
            position: 'is-bottom',
            type: 'is-danger'
          })
        })
    },
    updateMember: (id, names, email, phoneNumber, allergies, password, group, memberSetPassword) => {
      var member = {
        id,
        names,
        email,
        phoneNumber,
        allergies,
        password,
        group,
        memberSetPassword
      }
      axios({
        method: 'put',
        url: `/api/members/${id}`,
        data: member,
        headers: {
          Authorization: `Bearer ${localStorage.getItem('authToken')}`
        }})
        .then(resp => {
          console.log('Add successful', resp)
          Toast.open({
            message: 'Flatmate updated successfully',
            position: 'is-bottom',
            type: 'is-success'
          })
          /*
          setTimeout(() => {
            location.reload()
          }, 1000)
          */
        })
        .catch(err => {
          console.log('Add failed', err)
          Toast.open({
            message: 'Failed to update flatmate',
            position: 'is-bottom',
            type: 'is-danger'
          })
        })
    },
    deleteMember: (id) => {
      console.log('Attempting to remove member')
      Dialog.confirm({
        message: 'Are you sure you want to remove this flatmember? (this cannot be undone)',
        confirmText: 'Remove Flatmember',
        type: 'is-danger',
        onConfirm: () => {
          axios.delete(`/api/members/${id}`,
            {
              headers: {
                Authorization: `Bearer ${localStorage.getItem('authToken')}`
              }
            })
            .then(resp => {
              console.log('Remove successful', resp)
              Toast.open({
                message: 'Flatmate removed successfully',
                position: 'is-bottom',
                type: 'is-success'
              })
              setTimeout(() => {
                location.href = '#/admin/members'
              }, 1.2 * 1000)
            })
            .catch(err => {
              console.log('Remove failed', err)
              Toast.open({
                message: 'Failed to remove flatmate',
                position: 'is-bottom',
                type: 'is-danger'
              })
            })
        }
      })
    },
    returnNamesforID: (id, names) => {
      if (typeof id === 'undefined') {
        return null
      }
      return names
    },
    printGroupSelection: (group) => {
      console.log(group)
    }
  },
  components: {
    headerDisplay
  }
}
</script>

<style scoped>

</style>
