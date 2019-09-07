<template>
    <div>
        <div id="topbar" v-if="task && Object.keys(task).length">
            <h1>Task - {{ task.name }}</h1>
        </div>
        <div id="content">
            <div id="taskBody" v-if="task && Object.keys(task).length">
              <div id="meta">
                <p>
                  Importance: {{ task.importance }}/5 <br>
                  Description: {{ task.description }} <br>
                  Location: {{ task.location }}
                </p>
              </div>
            </div>
            <div id="entryform">
              <h2>Add an entry</h2>
              <form @submit.prevent="handleSubmit">
                <h3>I am</h3>
                <select name="member" id="iam_selector" v-model="form.names">
                  <option disabled value="">Please select your name</option>
                  <option disabled value="">--------</option>
                  <option v-bind:value="`${member.names}`" v-for="member of members" v-bind:key="member">
                    {{ member.names }}
                  </option>
                </select>
                <br><br>
                <input type="hidden" name="taskName" v-model="form.taskName">
                <input type="submit" value="Add entry">
              </form>
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
    axios.get(`http://localhost:3000/task/${this.$route.query.task}`)
      .then(response => {
        this.task = response.data
        this.form.taskName = response.data.name

        return axios.get(`http://localhost:3000/members`)
      })
      .then(response => {
        this.members = response.data
      })
      .catch(err => {
        this.tasksGETerrors.push(err)
      })
  },
  methods: {
    handleSubmit () {
      axios.post(`http://localhost:3000/entry`, this.form)
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
