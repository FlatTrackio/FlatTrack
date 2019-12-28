<template>
  <div>
    <headerDisplay />
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
          <ul>
            <li>
              <a href="/">Home</a>
            </li>
            <li class="is-active">
              <a href="/tasks">Tasks</a>
            </li>
          </ul>
        </nav>
        <h1 class="title">Tasks</h1>
        <h2 class="subtitle">Get caught up with your tasks</h2>
        <b-field label="How often would you like to be notified about tasks?">
          <b-select
            placeholder="Twice a week"
            expanded
            rounded
            v-model="alertFrequency"
            @input="updateFrequency(alertFrequency)"
          >
            <option value="3">Three a week</option>
            <option value="2">Twice a week</option>
            <option value="1">Once a week</option>
            <option value="0">Never</option>
          </b-select>
        </b-field>
        <div id="tasks" v-if="entries && entries.length && tasks">
          <h2 class="subtitle">Here are your assigned tasks</h2>
          <div class="card-margin" v-for="entry of entries" v-bind:key="entry">
            <div class="card">
              <div class="card-content">
                <div class="media">
                  <div class="media-content">
                    <p class="title is-4">
                      {{ tasks[entry.taskID].name }}
                    </p>
                    <p class="subtitle is-6">@{{ tasks[entry.taskID].location }}</p>
                  </div>
                </div>
                <div class="content">{{ tasks[entry.taskID].description }}</div>
                <div class="field is-grouped is-grouped-multiline">
                  <div class="control">
                    <div class="tags has-addons">
                      <span class="tag">rotates</span>
                      <span class="tag is-info">{{ tasks[entry.taskID].rotation }}</span>
                    </div>
                  </div>
                  <div class="control">
                    <div class="tags has-addons">
                      <span class="tag">is</span>
                      <span :class="entry.status == 'completed' ? `tag is-success`: `tag is-danger`">{{ entry.status }}</span>
                    </div>
                  </div>
                </div>
                <div v-if="entry.status == 'completed'">
                  Marked completed at {{ formatTimestampWithTime(entry.timestamp) }}
                </div>
                <div v-else>
                  This task must be completed by {{ formatTimestamp(entry.completeBy) }}
                </div>
              </div>
              <footer class="card-footer">
                <a :href="`/tasks/t?task=${entry.id}`" class="card-footer-item">View</a>
              </footer>
            </div>
          </div>
          <div class="section">
            <p>{{ entries.length }} {{ entries.length === 1 ? 'task' : 'tasks' }}</p>
          </div>
        </div>
        <div id="tasks" v-if="!entries.length">
          <section class="section">
            <div class="card">
              <div class="card-content">
                <div class="media">
                  <div class="media-content">
                    <p class="title is-4">No tasks here</p>
                  </div>
                </div>
                <div class="content">
                  Nice! You're either all caught up, or no tasks have been assigned to you.
                  <br />
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
import axios from 'axios'
import { Service } from 'axios-middleware'
import headerDisplay from '@/components/header-display'
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
  name: 'tasks',
  data () {
    return {
      tasks: [],
      entries: [],
      entryCompletedTimestamp: '',
      alertFrequency: 2
    }
  },
  created () {
    // get list of assigned tasks
    axios
      .get(`/api/tasks`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('authToken')}`
        }
      })
      .then(resp => {
        var transformRespData = {}
        resp.data.map(item => {
          transformRespData[item.id] = item
          delete transformRespData[item.id].id
        })
        this.tasks = transformRespData

        // get a list of the entries for this week
        return axios
          .get(`/api/entry`, {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          })
      })
      .then(resp => {
        this.entries = resp.data

        // get profile information
        return axios
          .get(`/api/profile`, {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          })
      }).then(resp => {
        this.alertFrequency = resp.data.taskNotificationFrequency
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
    updateFrequency: (value) => {
      axios({
        method: 'post',
        url: `/api/profile`,
        data: {
          frequency: value
        },
        headers: {
          Authorization: `Bearer ${localStorage.getItem('authToken')}`
        }
      })
        .then(resp => {
          // console.log(resp)
        })
        .catch(err => {
          console.log(err)
          // TODO display UI error
        })
    },
    formatTimestampWithTime: (timestamp) => {
      return moment.unix(timestamp).format('DD/MM/YYYY HH:mm')
    },
    formatTimestamp: (timestamp) => {
      return moment.unix(timestamp).format('DD/MM/YYYY')
    }
  },
  components: {
    headerDisplay
  }
}
</script>

<style src='../../assets/style.css'></style>
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
