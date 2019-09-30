<template>
    <div>
      <headerDisplay admin="true"/>
      <div class="container">
        <section class="section">
          <nav class="breadcrumb has-arrow-separator" aria-label="breadcrumbs">
            <ul>
                <li><a href="/#/admin">Admin Home</a></li>
                <li><a href="/#/admin/members">Manage flatmates</a></li>
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
                <b-checkbox v-if="!id" v-model="memberSetPassword" native-value="true">Allow the new member to set the password?<br/><br/></b-checkbox>
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
                        name="group">
                        Admin - can access the admin pages
                    </b-radio>
                  </div>
                  <div class="field">
                    <b-radio v-model="group"
                        native-value="approver"
                        name="group">
                        Approver - can view and approve tasks
                    </b-radio>
                  </div>
                </div>
                <br>
                <div v-if="returnNamesforID(id, names)">
                  <b-button type="is-success">Update</b-button>
                  <b-button type="is-warning">Disable</b-button>
                  <b-button type="is-danger" @click="deleteMember(id)">Delete</b-button>
                </div>
                <div v-else>
                  <b-button type="is-success" native-type="submit" @click="addNewMember(names, email, allergies, password, group, memberSetPassword)">Add new flatmate</b-button>
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
import { ToastProgrammatic as Toast } from 'buefy'
import headerDisplay from '../common/header-display'

export default {
  name: 'Admin home',
  data () {
    return {
      id: this.$route.query.id,
      names: '',
      email: '',
      allergies: '',
      password: '',
      memberSetPassword: true,
      group: 'flatmember',
      pageLocation: location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : ''),
      member: {
      },
      pageErrors: []
    }
  },
  created () {
    var id = this.$route.query.id
    axios.get(`/api/members/${id}`,
      {
        headers: {
          Authorization: `Bearer ${sessionStorage.getItem('authToken')}`
        }
      })
      .then(response => {
        var member = response.data
        this.names = response.data.names
        this.email = member.email
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
    addNewMember: (names, email, allergies, password, group, memberSetPassword) => {
      var member = {
        names,
        email,
        allergies,
        password,
        group,
        memberSetPassword
      }
      console.log(JSON.stringify(member, null, 4))
      axios({
        method: 'post',
        url: `/api/members`,
        data: member,
        headers: {
          Authorization: `Bearer ${sessionStorage.getItem('authToken')}`
        }}).then(response => {
        console.log('Add successful', response)
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
    deleteMember: (id) => {
      console.log('Attempting to remove member')
      axios.delete(`/api/members/${id}`, {data: {id: id}},
        {
          headers: {
            Authorization: `Bearer ${sessionStorage.getItem('authToken')}`
          }
        })
        .then(response => {
          console.log('Remove successful', response)
          Toast.open({
            message: 'Flatmate removed successfully',
            position: 'is-bottom',
            type: 'is-success'
          })
          location.href = '#/admin/members?reload=1'
        })
        .catch(err => {
          console.log('Remove failed', err)
          Toast.open({
            message: 'Failed to remove flatmate',
            position: 'is-bottom',
            type: 'is-danger'
          })
        })
    },
    returnNamesforID: (id, names) => {
      if (typeof id !== 'undefined') {
        return names
      } else {
        return null
      }
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
