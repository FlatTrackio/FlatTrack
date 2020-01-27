<template>
    <div>
        <headerDisplay admin="true"/>
        <div class="container">
            <section class="section">
                <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
                    <ul>
                        <li><a href="/admin">Admin Home</a></li>
                        <li class="is-active"><a href="/admin/entries">Manage Assigned Tasks and Entries</a></li>
                    </ul>
                </nav>
                <h1 class="title is-2">Assigned Tasks and Entries</h1>
                <p class="subtitle is-4">Display all active assigned entries</p>
                <div id="entries" v-if="entries && entries.length">
                  <div class="card-margin" v-for="entry of entries" v-bind:key="entry">
                      <div class="card">
                        <div class="card-content">
                            <div class="media">
                              <div class="media-content">
                                  <p class="title is-4">
                                    {{ tasks[entry.taskID].name }}
                                  </p>
                                  <div class="control">
                                    <div class="tags has-addons">
                                      <span class="tag">is</span>
                                      <span :class="entry.status == 'completed' ? `tag is-success`: `tag is-danger`">{{ entry.status }}</span>
                                    </div>
                                  </div>
                                  <div class="control">
                                    <div class="tags has-addons">
                                      <span class="tag">assigned to</span>
                                      <span class="tag is-success">{{ members[entry.member].names }}</span>
                                    </div>
                                  </div>
                              </div>
                            </div>
                            <div class="content">
                              <p>{{ tasks[entry.taskID].description }}</p>
                              <br>
                              <p>Assigned at {{ formatTimestampWithTime(entry.timestampAssign) }}</p>
                              <br>
                              <div v-if="entry.status === 'completed'">
                                <p>Completed at {{ formatTimestampWithTime(entry.timestamp) }}</p>
                              </div>
                              <div v-else>
                                <p>Must be completed by {{ formatTimestamp(entry.completeBy) }}</p>
                              </div>
                            </div>
                        </div>
                      </div>
                  </div>
                  <div class="section">
                    <p>{{ entries.length }} {{ entries.length === 1 ? 'entry' : 'entries' }}</p>
                  </div>
                </div>
                <div id="entries" v-if="!entries.length">
                  <section class="section">
                    <div class="card">
                    <div class="card-content">
                        <div class="media">
                        <div class="media-content">
                            <p class="title is-4">No entries active</p>
                        </div>
                        </div>
                        <div class="content">
                            No tasks have been assigned to anyone
                            <br />
                        </div>
                    </div>
                </div>
            </section>
        </div>
    </div>
</template>

<script>
import axios from 'axios'
import moment from 'moment'
import headerDisplay from '@/components/header-display'

export default {
  name: 'manage-entries',
  data () {
    return {
      tasks: [],
      entries: [],
      members: []
    }
  },
  created () {
    axios
      .get(`/api/admin/tasks`, {
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
          .get(`/api/admin/entry`, {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          })
      })
      .then(resp => {
        this.entries = resp.data

        return axios
          .get(`/api/admin/members`, {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          })
      })
      .then(resp => {
        var transformRespData = {}
        resp.data.map(item => {
          transformRespData[item.id] = item
          delete transformRespData[item.id].id
        })
        this.members = transformRespData

        console.log(this.entries)
        console.log(this.tasks)
        console.log(this.members)
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

<style scoped>

</style>
