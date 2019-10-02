<template>
  <div>
    <headerDisplay/>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb has-arrow-separator" aria-label="breadcrumbs">
            <ul>
            <li><a href="/#/">Home</a></li>
            <li><a href="/#/members">Flatmates</a></li>
            <li class="is-active"><a href="/#/members">{{ member.names }}</a></li>
            </ul>
        </nav>
        <h1 class="title">{{ member.names }}</h1>
        <div class="card-margin">
          <div class="card">
            <div class="card-content">
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
          </div>
        </div>
        <div v-if="!Object.keys(member) && !Object.keys(member).length">
          <div class="card">
            <div class="card-content">
              <div class="media">
                <div class="media-content">
                  <p class="title is-4">Can't find flatmate</p>
                </div>
              </div>
              <div class="content">
                Hmm, there appears to have been an issue.<br/>
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
import headerDisplay from '../common/header-display'

export default {
  name: 'Shopping List',
  data () {
    return {
      id: this.$route.query.id,
      member: {
      },
      pageLocation: location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '')
    }
  },
  created () {
    var id = this.$route.query.id
    axios.get(`/api/members/${id}`,
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('authToken')}`
        }
      })
      .then(response => {
        this.member = response.data
        this.member.password = null
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
