<template>
    <div>
      <headerDisplay admin="true"/>
      <div class="container">
        <section class="section">
          <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
                <li><a href="/admin">Admin Home</a></li>
                <li class="is-active"><a href="/admin/admin-manage-members">Manage flatmates</a></li>
            </ul>
          </nav>
          <h1 class="title is-2">Manage Flatmates</h1>
          <p class="subtitle is-4">Add, remove, and update Flatmates</p>
          <b-button type="is-light" tag="a" href="/admin/members/u" rounded>Add new flatmate</b-button>
          <br><br>
          <div v-if="members && members.length">
            <div class="card-margin" v-for="member of members" v-bind:key="member">
              <div class="card">
                <div class="card-content">
                  <div class="media">
                    <div class="media-content">
                      <p class="title is-4">{{ member.names }}</p>
                    </div>
                  </div>
                  <div class="content">
                    <div v-if="member.phoneNumber">
                      Phone: <a :href="`tel:${member.phoneNumber}`">{{ member.phoneNumber }}</a><br/>
                    </div>
                    <div v-if="member.email">
                      Email: <a :href="`mailto:${member.email}`">{{ member.email }}</a><br/>
                    </div>
                    <div v-if="member.allergies">
                      Allergies: {{ member.allergies }}<br/>
                    </div>
                  </div>
                </div>
                <footer class="card-footer">
                  <a :href="`/admin/members/u?id=${member.id}`" class="card-footer-item">View</a>
                </footer>
              </div>
            </div>
            <div class="section">
              <p>{{ members.length }} flatmembers</p>
            </div>
          </div>
          <div v-if="!members && !members.length">
            <section class="section">
              <div class="card">
                <div class="card-content">
                  <div class="media">
                    <div class="media-content">
                      <p class="title is-4">No flatmates added</p>
                    </div>
                  </div>
                  <div class="content">
                    Hmm, it appears that you don't have an flatmates added.<br/>
                    Press the add button to start.
                  </div>
                </div>
              </div>
            </section>
          </div>
        </section>
      </div>
    </div>
</template>

<script>
import axios from 'axios'
import { Service } from 'axios-middleware'
import headerDisplay from '@/components/header-display'

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
      members: [],
      pageErrors: []
    }
  },
  created () {
    axios({
      method: 'get',
      url: `/api/admin/members`,
      headers: {
        Authorization: `Bearer ${localStorage.getItem('authToken')}`
      }}).then(resp => {
      this.members = resp.data
    }).catch(err => {
      this.pageErrors.push(err)
    })
  },
  components: {
    headerDisplay
  }
}
</script>

<style scoped>

</style>
