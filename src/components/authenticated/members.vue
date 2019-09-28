<template>
  <div>
    <headerDisplay/>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb has-arrow-separator" aria-label="breadcrumbs">
            <ul>
            <li><a href="/#/">Home</a></li>
            <li class="is-active"><a href="/#/members">Flatmates</a></li>
            </ul>
        </nav>
        <h1 class="title">Flatmates</h1>
        <p>These are your flatmates, make sure to get to know them {{ emojiSmile }}</p>
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
                <a :href="`${pageLocation}/#/members/u?id=${member.id}`" class="card-footer-item">View</a>
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
import headerDisplay from '../common/header-display'
import emoji from 'node-emoji'

export default {
  name: 'Shopping List',
  data () {
    return {
      members: [],
      pageLocation: location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : ''),
      emojiSmile: emoji.get('smile')
    }
  },
  created () {
    axios.get(`/api/members`,
      {
        headers: {
          Authorization: `Bearer ${sessionStorage.getItem('authToken')}`
        }
      })
      .then(response => {
        this.members = response.data
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
