<template>
    <div>
      <headerDisplay/>
      <div class="container">
        <section class="section">
          <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><a href="/">Home</a></li>
              <li><a href="/tasks">Tasks</a></li>
              <li class="is-active"><a href="/#" aria-current="page">{{ task.name }}</a></li>
            </ul>
          </nav>
          <h1 class="title">{{ task.name }}</h1>
          <div class="card">
            <div class="card-content">
              <div class="media">
                <div class="media-content">
                  <p class="subtitle is-6">@{{ task.location }}</p>
                </div>
              </div>
              <div class="content">
                <div class="content">{{ task.description }}</div>
                <div class="field is-grouped is-grouped-multiline">
                  <div class="control">
                    <div class="tags has-addons">
                      <span class="tag">rotates</span>
                      <span class="tag is-info">{{ task.rotation }}</span>
                    </div>
                  </div>
                  <div class="control">
                    <div class="tags has-addons">
                      <span class="tag">is</span>
                      <span :class="entry.status == 'completed' ? `tag is-success`: `tag is-danger`">{{ entry.status }}</span>
                    </div>
                  </div>
              </div>
            </div>
            <div v-if="entry.status == 'uncompleted'">
              This task must be completed by {{ entryCompleteBy }}
              <br/>
            </div>
            <footer class="card-footer" v-if="entry.status == 'uncompleted'">
              <b-button @click="markAsCompleted(entry.id)" outlined tag="a" type="is-success" class="card-footer-item">I've completed this task</b-button>
            </footer>
            <div v-else>
              Marked completed at {{ entryCompletedTimestamp }}
            </div>
          </div>
        </section>
      </div>
    </div>
</template>

<script>
import axios from 'axios'
import { Service } from 'axios-middleware'
import headerDisplay from '../common/header-display'
import { NotificationProgrammatic as Notification } from 'buefy'
import moment from 'moment'

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
      entry: {},
      entryCompletedTimestamp: '',
      entryCompleteBy: '',
      tasksGETerrors: [],
      form: {
        names: '',
        taskName: ''
      }
    }
  },
  created () {
    axios
      .get(`/api/entry/${this.$route.query.task}`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('authToken')}`
        }
      })
      .then(resp => {
        this.entry = resp.data
        this.entryCompleteBy = moment.unix(this.entry.completeBy).format('DD/MM/YYYY')
        if (this.entry.timestamp) {
          this.entryCompletedTimestamp = moment.unix(this.entry.timestamp).format('DD/MM/YYYY HH:mm')
        }

        // get a list of the entries for this week
        return axios
          .get(`/api/task/${resp.data.taskID}`, {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          })
      })
      .then(resp => {
        this.task = resp.data
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
    markAsCompleted (id) {
      axios.put(`/api/entry/${id}`, {status: 'completed'},
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('authToken')}`
          }
        })
        .then(resp => {
          console.log(resp)
          Notification.open({
            duration: 5000,
            message: `Sucessully marked as completed`,
            position: 'is-bottom-right',
            hasIcon: true
          })
          setTimeout(() => {
            location.reload()
          }, 1000)
        })
        .catch(err => {
          Notification.open({
            duration: 5000,
            message: `An error has occured: ${err}`,
            position: 'is-bottom-right',
            type: 'is-danger',
            hasIcon: true
          })
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
