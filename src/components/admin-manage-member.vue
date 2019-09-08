<template>
    <div>
        <section class="hero is-warning">
            <div class="hero-body">
                <p class="title">
                    FlatTrack (admin)
                </p>
                <p class="subtitle">
                    {{ deploymentName }}
                </p>
            </div>
        </section>
        <div class="container">
          <section class="section">
            <nav class="breadcrumb has-arrow-separator" aria-label="breadcrumbs">
              <ul>
                  <li><a href="/#/admin">Admin Home</a></li>
                  <li><a href="/#/admin/members">Manage flatmates</a></li>
                  <li class="is-active"><a>{{ names || 'Add a flatmate' }}</a></li>
              </ul>
            </nav>
            <h1 class="title">{{ names || 'Add a flatmate' }}</h1>
            <div class="card">
              <div class="card-content">
                <div class="content">
                  <b-field label="Fullname" v-if="!names">
                      <b-input placeholder="xxxxx xxxxxxxx" v-model="member.names" maxlength="30" rounded required></b-input>
                  </b-field>
                  <b-field label="Phone Number">
                      <b-input placeholder="xx xxx xxxx" v-model="member.phoneNumber" maxlength="30" rounded></b-input>
                  </b-field>
                  <b-field label="Email">
                      <b-input placeholder="xxxxx@xxxxx.xxx" v-model="member.email" maxlength="30" rounded required></b-input>
                  </b-field>
                  <b-field label="Allergies">
                      <b-input placeholder="xx, xxxx, xx" v-model="member.allergies" maxlength="30" rounded></b-input>
                  </b-field>
                  <b-field label="Password">
                      <b-input type="password" placeholder="xxxxxxxxxxxxx" v-model="member.password" maxlength="30" rounded :required="names" password-reveal></b-input>
                  </b-field>
                  <div v-if="names">
                    <b-button type="is-success">Update</b-button>
                    <b-button type="is-warning">Disable</b-button>
                    <b-button type="is-danger">Delete</b-button> 
                  </div>
                  <div v-if="!names">
                    <b-button type="is-success" @click="addNewMember">Add new flatmate</b-button>
                  </div>
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

export default {
  name: 'Admin home',
  data () {
    return {
      deploymentName: 'Keep track of your flat',
      id: this.$route.query.id,
      names: null,
      pageLocation: location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : ''),
      member: {}
    }
  },
  created () {
    if (this.$route.query.id != null) axios.get(`${location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '')}/api/members/${this.$route.query.id}`)
      .then(response => {
        this.member = response.data
        this.names = this.member.names
        this.member.password = null
      })
      .catch(err => {
        this.pageErrors.push(err)
      })
  },
  methods: {
    addNewMember: () => {
      axios.post(`/api/members`, this.member)
      .then(response => {
        console.log('Add successful', response)
        Toast.open({
          message: 'Flatmate added successfully',
          position: 'is-bottom',
          type: 'is-success'
        })
      })
      .catch(err => {
        console.log('Add failed', err)
        Toast.open({
          message: 'Failed to add flatmate',
          position: 'is-bottom',
          type: 'is-danger'
        })
        this.pageErrors.push(err)
      })
    }
  }
}
</script>

<style scoped>

</style>