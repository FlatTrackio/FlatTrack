<template>
    <div>
        <headerDisplay admin="true"/>
        <div class="container">
            <section class="section">
                <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
                    <ul>
                        <li><a href="/admin">Admin Home</a></li>
                        <li class="is-active"><a href="/admin/tasks">Manage Tasks</a></li>
                    </ul>
                </nav>
                <h1 class="title is-2">Manage Tasks</h1>
                <p class="subtitle is-4">Set up the tasks which need to be done</p>
                <b-button type="is-light" tag="a" href="/admin/tasks/t" rounded>Add a new task</b-button>
                <br><br>
                <div id="tasks" v-if="tasks && tasks.length">
                    <div class="card-margin" v-for="task of tasks" v-bind:key="task">
                        <div class="card">
                        <div class="card-content">
                            <div class="media">
                            <div class="media-content">
                                <p class="title is-4">
                                {{ task.name }}
                                <span class="tag is-info">rotates {{ task.rotation }}</span>
                                </p>
                                <p class="subtitle is-6">@{{ task.location }}</p>
                            </div>
                            </div>
                            <div class="content">{{ task.description }}</div>
                        </div>
                        <footer class="card-footer">
                            <a :href="`/admin/tasks/t?id=${task.id}`" class="card-footer-item">Edit</a>
                        </footer>
                        </div>
                    </div>
                    <div class="section">
                      <p>{{ tasks.length }} {{ tasks.length === 1 ? 'task' : 'tasks' }}</p>
                    </div>
                </div>
                <div id="tasks" v-if="!tasks.length">
                <section class="section">
                    <div class="card">
                    <div class="card-content">
                        <div class="media">
                        <div class="media-content">
                            <p class="title is-4">No tasks added</p>
                        </div>
                        </div>
                        <div class="content">
                            Go ahead and add some, so everyone can keep track of their tasks
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
import headerDisplay from '@/components/header-display'

export default {
  name: 'manage-tasks',
  data () {
    return {
      tasks: []
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
        this.tasks = resp.data
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
  components: {
    headerDisplay
  }
}
</script>

<style scoped>

</style>
