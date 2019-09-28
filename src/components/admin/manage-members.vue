<template>
    <div>
      <headerDisplay admin="true"/>
      <div class="container">
        <section class="section">
          <nav class="breadcrumb has-arrow-separator" aria-label="breadcrumbs">
            <ul>
                <li><a href="/#/admin">Admin Home</a></li>
                <li class="is-active"><a href="/#/admin/admin-manage-members">Manage flatmates</a></li>
            </ul>
          </nav>
          <h1 class="title">Manage Flatmates</h1>
          <h2 class="subtitle">Add, remove, and update Flatmates</h2>
          <b-button type="is-light" @click="addNewFlatmate" rounded>Add new flatmate</b-button>
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
import headerDisplay from '../common/header-display'

export default {
  name: 'Admin home',
  data () {
    return {
      pageLocation: location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : ''),
      members: []
    }
  },
  methods: {
    addNewFlatmate: () => {
      location.href = `/#/admin/members/u`
    }
  },
  components: {
    headerDisplay
  }
}
</script>

<style scoped>

</style>
