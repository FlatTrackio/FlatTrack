<template>
  <div>
    <headerDisplay/>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb has-arrow-separator" aria-label="breadcrumbs">
            <ul>
            <li><a href="/#/">Home</a></li>
            <li class="is-active"><a href="/#/flatmates">Flatmates</a></li>
            </ul>
        </nav>
        <h1 class="title">Flatmates</h1>
        <p>These are your flatmates, make sure to get to know them 😃</p>
      </section>
      <section class="section">
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
                    Phone: {{ member.phoneNumber }}<br/>
                  </div>
                  <div v-if="member.email">
                    Email: {{ member.email }}<br/>
                  </div>
                  <div v-if="member.allergies">
                    Allergies: {{ member.allergies }}<br/>
                  </div>
                </div>
              </div>
              <footer class="card-footer">
                <a :href="`${pageLocation}/#/admin/members/u?id=${member.id}`" class="card-footer-item">View</a>
              </footer>
            </div>
          </div>
        </div>
        <div v-if="!members && !members.length">
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
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import headerDisplay from './header-display'

export default {
  name: 'Shopping List',
  data () {
    return {
      deploymentName: 'Keep track of your flat',
      members: [],
      pageLocation: location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '')
    }
  },
  created () {
    axios.get(`/api/members`)
      .then(response => {
        this.members = response.data

        return axios.get(`/api/settings/deploymentName`)
      })
      .then(response => {
        this.deploymentName = response.data.value
      })
      .catch(err => {
        this.$buefy.notification.open({
          duration: 5000,
          message: `An error has occured: ${err}`,
          position: 'is-bottom-right',
          type: 'is-danger',
          hasIcon: true
        })
      })
  },
  methods: {
  },
  components: {
    headerDisplay
  }
}
</script>

<style scoped>

</style>
