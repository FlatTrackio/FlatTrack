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
                  <p class="title is-4">
                    {{ task.name }} <span class="tag is-info">rotates {{ task.rotation }}</span>
                  </p>
                  <p class="subtitle is-6">@{{ task.location }}</p>
                </div>
              </div>
              <div class="content">
                {{ task.description }}
              </div>
            </div>
          </div>
          <br>
          <b-button @click="markAsCompleted" rounded outlined type="is-success">I've completed this task</b-button>
        </section>
      </div>
    </div>
</template>

<script>
import axios from 'axios'
import { Service } from 'axios-middleware'
import headerDisplay from '../common/header-display'

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
    axios.get(`/api/task/${this.$route.query.task}`,
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('authToken')}`
        }
      })
      .then(resp => {
        this.task = resp.data
        this.form.taskName = resp.data.name

        return axios.get(`/api/members`,
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          })
      })
      .then(resp => {
        this.members = resp.data
      })
      .catch(err => {
        this.tasksGETerrors.push(err)
      })
  },
  methods: {
    markAsCompleted () {
      axios.post(`/api/entry/${this.task.id}`, {},
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('authToken')}`
          }
        })
        .then(resp => {
          console.log(resp)
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
