<template>
  <div>
    <headerDisplay/>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb has-arrow-separator" aria-label="breadcrumbs">
          <ul>
            <li><a href="/#/">Home</a></li>
            <li class="is-active"><a href="/#/tasks">Tasks</a></li>
          </ul>
        </nav>
        <h1 class="title">Tasks</h1>
        <h2 class="subtitle">Get caught up with your tasks</h2>
        <div id="tasks" v-if="tasks && tasks.length">
          <b-field label="How often would you like to be notified about tasks?">
            <b-select
                placeholder="Medium"
                expanded
                rounded
                v-model="alertFrequency">
                <option value="3">Three a week</option>
                <option value="2">Twice a week</option>
                <option value="1">Once a week</option>
                <option value="0">Never</option>
            </b-select>
          </b-field>
          <h2 class="subtitle">Here are your assigned tasks</h2>
          <div class="card-margin" v-for="task of tasks" v-bind:key="task">
            <div class="card">
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
                <a :href="`${pageLocation}/#/tasks/view?task=${task.id}`" class="card-footer-item">View</a>
              </footer>
            </div>
          </div>
      </div>
      <div id="tasks" v-if="!tasks && !tasks.length">
        <section class="section">
          <div class="card">
            <div class="card-content">
              <div class="media">
                <div class="media-content">
                  <p class="title is-4">No tasks here</p>
                </div>
              </div>
              <div class="content">
                Nice! You're either all caught up, or no tasks have been assigned to you.<br/>
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
  name: 'tasks',
  data () {
    return {
      deploymentName: 'Keep track of your flat',
      tasks: [],
      members: [],
      pageLocation: location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : ''),
      alertFrequency: 2
    }
  },
  created () {
    axios.get(`/api/tasks`)
      .then(response => {
        this.tasks = response.data

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
