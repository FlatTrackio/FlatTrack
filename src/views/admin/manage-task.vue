<template>
    <div>
      <headerDisplay admin="true"/>
      <div class="container">
        <section class="section">
          <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><a href="/admin">Admin Home</a></li>
              <li><a href="/admin/tasks">Manage Tasks</a></li>
              <li class="is-active"><a>{{ returnNameforID(id, name) || 'Add a task' }}</a></li>
            </ul>
          </nav>
          <h1 class="title">{{ returnNameforID(id, name) || 'Add a task' }}</h1>
          <div class="card">
            <div class="card-content">
              <div class="content">
                <b-field label="Name" v-if="!id">
                  <b-input placeholder="What should this task be called?" v-model="name" maxlength="30" rounded required></b-input>
                </b-field>
                <b-field label="Description">
                  <b-input placeholder="Describe this task. i.e: Examples" v-model="description" maxlength="100" rounded required></b-input>
                </b-field>
                <b-field label="Location">
                  <b-input placeholder="Where is this to be found?" v-model="location" maxlength="100" rounded required></b-input>
                </b-field>
                <b-field label="Rotation">
                  <b-select placeholder="How often this task is rotated?" expanded rounded v-model="rotation" required>
                    <option value="monthly">Monthly</option>
                    <option value="weekly">Weekly</option>
                    <option value="daily">Daily</option>
                    <option value="never">Never</option>
                  </b-select>
                </b-field>
                <br>
                <b-field label="Frequency" v-if="rotation == 'never'">
                  <b-select placeholder="How often this task is done?" expanded rounded v-model="frequency" :required="rotation == 'never'">
                    <option value="monthly">Monthly</option>
                    <option value="weekly">Weekly</option>
                    <option value="daily">Daily</option>
                  </b-select>
                </b-field>
                <br>
                <b-field label="Assignee" v-if="rotation == 'never'">
                  <b-select placeholder="Who should be assigned to this task?" expanded rounded v-model="assignee" :required="rotation == 'never'">
                    <option v-for="member of members" v-bind:key="member" :value="member.id">{{ member.names }}</option>
                  </b-select>
                </b-field>
               </div>
              <div v-if="returnNameforID(id, name)">
                <b-button type="is-success" native-type="submit" @click="updateTask(id, name, description, location, disabled, assignee, rotation, frequency)">Update</b-button>
                <b-button type="is-warning">Disable</b-button>
                <b-button type="is-danger" @click="deleteTask(id)">Delete</b-button>
              </div>
              <div v-else>
                <b-button type="is-success" native-type="submit" @click="addNewTask(name, description, location, disabled, assignee, rotation, frequency)">Add new task</b-button>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>
</template>

<script>
import axios from 'axios'
import { Service } from 'axios-middleware'
import { ToastProgrammatic as Toast, DialogProgrammatic as Dialog } from 'buefy'
import headerDisplay from '@/components/header-display'

const service = new Service(axios)
service.register({
  onResponse (response) {
    if (response.status === 403) {
      localStorage.removeItem('authToken')
      window.location.href = '/'
    }
    return response
  }
})
export default {
  name: 'Admin home',
  data () {
    return {
      id: this.$route.query.id,
      name: '',
      description: '',
      location: '',
      disabled: false,
      assignee: null,
      rotation: null,
      frequency: null,
      members: [],
      pageErrors: []
    }
  },
  created () {
    var id = this.$route.query.id
    axios({
      method: 'get',
      url: `/api/admin/members`,
      headers: {
        Authorization: `Bearer ${localStorage.getItem('authToken')}`
      }})
      .then(resp => {
        this.members = resp.data

        if (!id) {
          return
        }
        return axios.get(`/api/admin/task/${id}`,
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          })
      })
      .then(resp => {
        console.log(resp.data)
        this.name = resp.data.name
        this.description = resp.data.description
        this.location = resp.data.location
        this.disabled = resp.data.disabled
        this.assignee = resp.data.assignee
        this.rotation = resp.data.rotation
        this.frequency = resp.data.frequency
      })
      .catch(err => {
        this.pageErrors = [...this.pageErrors, err]
      })
  },
  methods: {
    addNewTask: (name, description, location, disabled, assignee, rotation, frequency) => {
      var task = {
        name,
        description,
        location,
        disabled,
        assignee,
        rotation,
        frequency
      }
      axios({
        method: 'post',
        url: `/api/admin/tasks`,
        data: task,
        headers: {
          Authorization: `Bearer ${localStorage.getItem('authToken')}`
        }}).then(resp => {
        Toast.open({
          message: 'Task added successfully',
          position: 'is-bottom',
          type: 'is-success'
        })
        setTimeout(() => {
          window.window.location.href = 'admin/tasks'
        }, 1000)
      })
        .catch(err => {
          Toast.open({
            message: `Failed to add task: ${err}`,
            position: 'is-bottom',
            type: 'is-danger'
          })
        })
    },
    updateTask: (id, name, description, location, disabled, assignee, rotation) => {
      var task = {
        name,
        description,
        location,
        disabled,
        assignee,
        rotation
      }
      axios({
        method: 'put',
        url: `/api/admin/task/${id}`,
        data: task,
        headers: {
          Authorization: `Bearer ${localStorage.getItem('authToken')}`
        }})
        .then(resp => {
          Toast.open({
            message: 'Task updated successfully',
            position: 'is-bottom',
            type: 'is-success'
          })
          setTimeout(() => {
            window.location.reload()
          }, 1000)
        })
        .catch(err => {
          Toast.open({
            message: `Failed to update task: ${err}`,
            position: 'is-bottom',
            type: 'is-danger'
          })
        })
    },
    deleteTask: (id) => {
      Dialog.confirm({
        message: 'Are you sure you want to remove this task? (this cannot be undone)',
        confirmText: 'Remove task',
        type: 'is-danger',
        onConfirm: () => {
          axios.delete(`/api/admin/task/${id}`,
            {
              headers: {
                Authorization: `Bearer ${localStorage.getItem('authToken')}`
              }
            })
            .then(resp => {
              Toast.open({
                message: 'Task removed successfully',
                position: 'is-bottom',
                type: 'is-success'
              })
              setTimeout(() => {
                window.window.location.href = 'admin/tasks'
              }, 1.2 * 1000)
            })
            .catch(err => {
              Toast.open({
                message: `Failed to remove task: ${err}`,
                position: 'is-bottom',
                type: 'is-danger'
              })
            })
        }
      })
    },
    returnNameforID: (id, name) => {
      if (typeof id === 'undefined') {
        return null
      }
      return name
    }
  },
  components: {
    headerDisplay
  }
}
</script>

<style scoped>

</style>
