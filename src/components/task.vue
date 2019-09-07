<template>
    <div>
      <section class="hero is-dark">
        <div class="hero-body">
          <p class="title">
            Task - {{ task.name }}
          </p>
        </div>
      </section>
      <div class="container">

        <nav class="breadcrumb has-arrow-separator" aria-label="breadcrumbs">
          <ul>
            <li><a href="/#/">Home</a></li>
            <li><a href="/#/tasks">Tasks</a></li>
            <li class="is-active"><a href="/#" aria-current="page">{{ task.name }}</a></li>
          </ul>
        </nav>
        <div id="content">
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
            <section>
                <b-button @click="markAsCompleted">I've completed this task</b-button>
            </section>
          </div>
        </div>
        <div id="tasksGETerrors" v-if="tasksGETerrors && tasksGETerrors.length">
            <h3>Uh oh! Something's gone wrong with this page:</h3>
            <ul>
                <li v-for="tasksGETerror of tasksGETerrors" v-bind:key="tasksGETerror">
                    {{ tasksGETerror.message }}
                </li>
            </ul>
        </div>
      </div>
    </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'Task',
  data () {
    return {
      id: this.$route.query.task,
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
    axios.get(`${location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '')}/api/task/${this.$route.query.task}`)
      .then(response => {
        this.task = response.data
        this.form.taskName = response.data.name
        console.log(this.task)

        return axios.get(`${location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '')}/api/members`)
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
      axios.put(`${location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '')}/api/entry/${this.task.id}`, {})
        .then(response => {
          console.log(response)
        })
        .catch(err => {
          this.tasksGETerrors.push(err)
        })
    }
  }
}
</script>

<style src="../assets/style.css"></style>
<style scoped>
</style>
