<template>
    <div>
      <headerDisplay/>
      <div class="container">
        <section class="section">
          <nav class="breadcrumb has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><a href="/#/">Home</a></li>
              <li><a href="/#/tasks">Tasks</a></li>
              <li class="is-active"><a href="/#" aria-current="page">{{ task.name }}</a></li>
            </ul>
          </nav>
        </section>
        <section class="section">
          <h1 class="title">{{ task.name }}</h1>
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
              <b-button @click="markAsCompleted">I've completed this task</b-button>
            </footer>
          </div>
        </section>
      </div>
    </div>
</template>

<script>
import axios from 'axios'
import headerDisplay from '../common/header-display'

export default {
  name: 'Task',
  data () {
    return {
      id: this.$route.query.task,
      pageLocation: location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : ''),
      task: {},
      members: [],
      tasksGETerrors: [],
      form: {
        names: '',
        taskName: ''
      }
    }
  },
  created () {
    axios.get(`/api/task/${this.$route.query.task}`)
      .then(response => {
        this.task = response.data
        this.form.taskName = response.data.name

        return axios.get(`/api/members`)
      })
      .then(response => {
        this.members = response.data
      })
      .catch(err => {
        this.tasksGETerrors.push(err)
      })
  },
  methods: {
    markAsCompleted () {
      axios.put(`/api/entry/${this.task.id}`, {})
        .then(response => {
          console.log(response)
        })
        .catch(err => {
          this.tasksGETerrors.push(err)
        })
    }
  },
  components: {
    headerDisplay
  }
}
</script>

<style src="../../assets/style.css"></style>
<style scoped>
</style>
