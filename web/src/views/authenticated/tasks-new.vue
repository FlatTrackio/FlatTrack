<!--
     This program is free software: you can redistribute it and/or modify
     it under the terms of the Affero GNU General Public License as published by
     the Free Software Foundation, either version 3 of the License, or
     (at your option) any later version.

     This program is distributed in the hope that it will be useful,
     but WITHOUT ANY WARRANTY; without even the implied warranty of
     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
     GNU General Public License for more details.

     You should have received a copy of the Affero GNU General Public License
     along with this program.  If not, see <https://www.gnu.org/licenses/>.
-->

<template>
  <div>
    <div class="container">
      <section class="section">
        <nav
          class="breadcrumb is-medium has-arrow-separator"
          aria-label="breadcrumbs"
        >
          <ul>
            <li>
              <router-link :to="{ name: 'Tasks' }">Tasks</router-link>
            </li>
            <li class="is-active">
              <router-link :to="{ name: 'New task' }">New task</router-link>
            </li>
          </ul>
        </nav>
        <div>
          <h1 class="title is-1">New task</h1>
          <p class="subtitle is-3">Schedule a new task</p>
          <b-field label="Name">
            <b-input
              type="text"
              v-model="name"
              maxlength="30"
              icon="textbox"
              size="is-medium"
              placeholder="Enter a title for this task"
              autofocus
              icon-right="close-circle"
              icon-right-clickable
              @icon-right-click="name = ''"
              @keyup.enter.native="PostNewTask"
              required
            >
            </b-input>
          </b-field>
          <b-field label="Notes (optional)">
            <b-input
              type="text"
              size="is-medium"
              v-model="notes"
              icon="text"
              placeholder="Enter extra information"
              @keyup.enter.native="PostNewTask"
              icon-right="close-circle"
              icon-right-clickable
              @icon-right-click="notes = ''"
              maxlength="100"
            >
            </b-input>
          </b-field>
          <b-field label="Assignee type">
            <b-select
              placeholder="How flatmates should be assigned to this task."
              v-model="assigneeType"
              icon="account-convert-outline"
              expanded
              size="is-medium"
              required
            >
              <option value="next" default>Next</option>
              <option value="random">Random</option>
              <option value="self">Myself</option>
            </b-select>
          </b-field>
          <b-field label="Frequency">
            <b-select
              placeholder="When the task should occur"
              v-model="frequency"
              icon="clock-start"
              expanded
              size="is-medium"
              required
            >
              <option value="once">Once</option>
              <option value="daily">Daily</option>
              <option value="weekly">Weekly</option>
              <option value="fortnightly">Fortnightly</option>
              <option value="monthly">Monthly</option>
            </b-select>
          </b-field>
          <b-field label="Select when it should first occur">
            <b-datetimepicker
              v-model="targetStartTimestampPicker"
              size="is-medium"
              placeholder="Click to select..."
              icon="calendar-today"
              :icon-right="targetStartTimestampPicker ? 'close-circle' : ''"
              icon-right-clickable
              @icon-right-click="
                () => {
                  targetStartTimestampPicker = null;
                }
              "
              first-day-of-week="1"
              datepicker="{ showWeekNumber: true }"
              timepicker="{ enableSeconds: false, hourFormat: undefined }"
              horizontal-time-picker
            >
            </b-datetimepicker>
          </b-field>
          <b-field label="Template task (optional)" v-if="tasks.length > 0">
            <b-select
              placeholder="Optionally select a task to base a new task off"
              v-model="taskTemplate"
              icon="content-copy"
              expanded
              size="is-medium"
            >
              <option value="">No template</option>
              <option disabled></option>
              <option v-for="task in tasks" :value="task.id" :key="task.id">
                {{ task.name }}
              </option>
            </b-select>
          </b-field>
          <div
            class="field"
            v-if="taskTemplate !== '' && typeof taskTemplate !== 'undefined'"
          >
            <label class="label"> Select items </label>
            <div class="field">
              <b-radio
                v-model="templateTaskItemSelector"
                size="is-medium"
                name="itemSelector"
                native-value="all"
              >
                All items
              </b-radio>
            </div>
            <div class="field">
              <b-radio
                v-model="templateTaskItemSelector"
                size="is-medium"
                name="itemSelector"
                native-value="unobtained"
              >
                Only from unobtained
              </b-radio>
            </div>
            <div class="field">
              <b-radio
                v-model="templateTaskItemSelector"
                size="is-medium"
                name="itemSelector"
                native-value="obtained"
              >
                Only from obtained
              </b-radio>
            </div>
          </div>
          <b-button
            icon-left="plus"
            type="is-success"
            size="is-medium"
            native-type="submit"
            expanded
            :loading="submitLoading"
            :disabled="submitLoading"
            @click="PostNewTask"
          >
            Create task
          </b-button>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/common/common'
import tasks from '@/requests/authenticated/tasks'

export default {
  name: 'tasks-new',
  data () {
    return {
      name: '',
      notes: '',
      assigneeType: 'next',
      targetStartTimestamp: 0,
      targetStartTimestampPicker: 0,
      frequency: 'monthly',
      taskTemplate: '',
      templateListItemSelector: 'all',
      tasks: []
    }
  },
  methods: {
    PostNewTask () {
      if (this.notes === '') {
        this.notes = undefined
      }
      this.submitLoading = true
      tasks
        .PostTask(
          this.name,
          this.notes,
          this.assigneeType,
          this.targetStartTimestamp,
          this.frequency,
          this.taskTemplate
        )
        .then((resp) => {
          this.submitLoading = false
          var task = resp.data.spec
          if (task.id !== '' || typeof task.id === 'undefined') {
            this.$router.push({
              name: 'Edit task',
              params: { id: task.id }
            })
          } else {
            common.DisplayFailureToast('Unable to find created task')
          }
        })
        .catch((err) => {
          this.submitLoading = false
          common.DisplayFailureToast(
            `Failed to create task - ${err.response.data.metadata.response}`
          )
        })
    }
  },
  async beforeMount () {
    tasks.GetTasks(undefined, 'templated').then((resp) => {
      this.tasks = resp.data.list || []
    })
    this.name = this.$route.query.name
  }
}
</script>
