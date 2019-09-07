<template>
  <div>
    <section class="hero is-dark">
      <div class="hero-body">
        <p class="title">
          FlatTrack - {{ deploymentName }}
        </p>
      </div>
    </section>
    <div class="container">
      <nav class="breadcrumb has-arrow-separator" aria-label="breadcrumbs">
        <ul>
          <li><a href="/#/">Home</a></li>
          <li class="is-active"><a href="/#/tasks">Tasks</a></li>
        </ul>
      </nav>
      <div id="tasks" v-if="tasks && tasks.length">
        <h1 class="title">Tasks</h1>
        <div class="card" v-for="task of tasks" v-bind:key="task">
          <div class="card-content">
            <div class="media">
              <div class="media-content">
                <p class="title is-4">{{ task.name }}</p>
                <p class="subtitle is-6">@{{ task.location }}</p>
              </div>
            </div>
            <div class="content">
              {{ task.description }}
            </div>
          </div>
          <footer class="card-footer">
            <a :href="`${pageLocation}/#/task?task=${task.id}`" class="card-footer-item">View</a>
          </footer>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'tasks',
  data () {
    return {
      deploymentName: '[Deployment name]',
      tasks: [],
      members: [],
      pageLocation: location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '')
    }
  },
  created () {
    axios.get(`${location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '')}/api/tasks`)
      .then(response => {
        this.tasks = response.data

        return axios.get(`${location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '')}/api/members`)
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
  }
}
</script>

<style src="../assets/style.css"></style>
<style scoped>
.taskItem .child {
  background-color: lightblue;
  padding-top: 10px;
  padding-left: 10px;
  padding-bottom: 10px;
  margin-top: 5px;
  margin-bottom: 5px;
}

.memberItem {
  background-color: lightblue;
  padding-top: 10px;
  padding-left: 10px;
  padding-bottom: 10px;
  margin-top: 5px;
  margin-bottom: 5px;
}
</style>
