<template>
  <div>
    <div id="topbar">
      <h1>FlatTrack - {{ deploymentName }}</h1>
    </div>
    <div id="content">
      <div id="tasks" v-if="tasks && tasks.length">
        <h3>Tasks</h3>
        <div class="taskItem" v-for="task of tasks" v-bind:key="task">
            <div class="child" v-bind:onclick="`window.location='http://localhost:8080/#/task?task=${task.id}'`">
              <h4>{{ task.name }} ({{ task.importance }}/5 importance)</h4>
              <p>
                Description: {{ task.description }} <br>
                Location: {{ task.location }}
              </p>
            </div>
        </div>
      </div>
      <div id="members" v-if="members && members.length">
        <h3>Flatmates</h3>
        <div class="memberItem" v-for="member of members" v-bind:key="member">
          <h4>{{ member.names }}</h4>
        </div>
      </div>
      <div id="tasksGETerrors" v-if="tasksGETerrors && tasksGETerrors.length">
        <h3>Uh oh! Something's gone wrong with this page:</h3>
        <ul>
          <li v-for="tasksGETerror of tasksGETerrors" v-bind:key="tasksGETerror">
            {{tasksGETerror.message}}
          </li>
        </ul>
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
      tasksGETerrors: []
    }
  },
  created () {
    axios.get(`${location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '')}/api/task`)
      .then(response => {
        this.tasks = response.data

        return axios.get(`${location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '')}/api/members`)
      })
      .then(response => {
        this.members = response.data
      })
      .catch(err => {
        this.tasksGETerrors.push(err)
      })
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
