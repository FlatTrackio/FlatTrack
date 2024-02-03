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
            <li><router-link :to="{ name: 'Apps' }">Apps</router-link></li>
            <li class="is-active">
              <router-link :to="{ name: 'Tasks' }">Tasks</router-link>
            </li>
          </ul>
        </nav>
        <h1 class="title is-1">Tasks</h1>
        <p class="subtitle is-3">Get flat chores done</p>
        <div>
          <b-tabs
            :position="deviceIsMobile ? 'is-centered' : ''"
            class="block is-marginless"
            v-model="taskDisplayState"
          >
            <b-tab-item icon="format-list-checks" label="All"></b-tab-item>
            <b-tab-item icon="playlist-remove" label="Uncompleted"></b-tab-item>
            <b-tab-item icon="playlist-check" label="Completed"></b-tab-item>
          </b-tabs>
          <label class="label">Search for tasks</label>
          <b-field>
            <b-input
              icon="magnify"
              size="is-medium"
              placeholder="Enter a task name"
              type="search"
              expanded
              v-model="taskSearch"
              ref="search"
            >
            </b-input>
            <p class="control">
              <b-select
                placeholder="Sort by"
                icon="sort"
                v-model="sortBy"
                size="is-medium"
                expanded
              >
                <option value="recentlyAdded">Recently Added</option>
                <option value="lastAdded">Last Added</option>
                <option value="recentlyUpdated">Recently Updated</option>
                <option value="lastUpdated">Last Updated</option>
                <option value="alphabeticalDescending">A-z</option>
                <option value="alphabeticalAscending">z-A</option>
              </b-select>
            </p>
          </b-field>
          <b-loading
            :is-full-page="false"
            :active.sync="pageLoading"
            :can-cancel="false"
          ></b-loading>
          <section>
            <div
              class="card pointer-cursor-on-hover"
              @click="
                $router.push({
                  name: 'New task',
                  query: { name: taskSearch || undefined },
                })
              "
            >
              <div class="card-content">
                <div class="media">
                  <div class="media-left">
                    <b-icon icon="clipboard-plus-outline" size="is-medium">
                    </b-icon>
                  </div>
                  <div class="media-content">
                    <p class="title is-4">Add a new task</p>
                  </div>
                  <div class="media-right">
                    <b-icon
                      icon="chevron-right"
                      size="is-medium"
                      type="is-midgray"
                    ></b-icon>
                  </div>
                </div>
              </div>
            </div>
          </section>
        </div>
        <floatingAddButton
          :routerLink="{
            name: 'New task',
            query: { name: taskSearch || undefined },
          }"
          v-if="displayFloatingAddButton"
        />
        <br />
        <div v-if="tasksFiltered.length > 0">
          <taskListCardView
            :task="task"
            :authors="authors"
            :tasks="tasks"
            :index="index"
            v-for="(task, index) in tasksFiltered"
            v-bind:key="task"
            :deviceIsMobile="deviceIsMobile"
            @tasks="
              (l) => {
                tasks = l;
              }
            "
          />
          <br />
          <p>{{ tasksFiltered.length }} task(s)</p>
        </div>
        <div v-else>
          <div class="card">
            <div class="card-content card-content-list">
              <div class="media">
                <div class="media-left">
                  <b-icon
                    icon="cart-remove"
                    size="is-medium"
                    type="is-midgray"
                  ></b-icon>
                </div>
                <div class="media-content">
                  <p
                    class="subtitle is-4"
                    v-if="
                      taskSearch === '' && tasks.length === 0 && !pageLoading
                    "
                  >
                    No tasks added yet.
                  </p>
                  <p
                    class="subtitle is-4"
                    v-else-if="
                      taskSearch === '' &&
                      taskDisplayState === 1 &&
                      tasks.length > 0 &&
                      !pageLoading
                    "
                  >
                    All tasks have been completed.
                  </p>
                  <p
                    class="subtitle is-4"
                    v-else-if="
                      taskSearch === '' &&
                      taskDisplayState === 2 &&
                      tasks.length > 0 &&
                      !pageLoading
                    "
                  >
                    No tasks have been completed yet.
                  </p>
                  <p
                    class="subtitle is-4"
                    v-else-if="taskSearch !== '' && !pageLoading"
                  >
                    No tasks found.
                  </p>
                  <p class="subtitle is-4" v-else-if="pageLoading">
                    Loading tasks...
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/common/common'
import tasks from '@/requests/authenticated/tasks'
import flatmates from '@/requests/authenticated/flatmates'

export default {
  name: 'flattrack-tasks',
  data () {
    return {
      displayFloatingAddButton: true,
      canUserAccountAdmin: false,
      notes: '',
      tasks: [],
      authors: {},
      taskDisplayState: 0,
      deviceIsMobile: false,
      taskSearch: '',
      pageLoading: true,
      sortBy: 'recentlyUpdated'
    }
  },
  components: {
    taskListCardView: () =>
      import('@/components/authenticated/task-list-card-view.vue'),
    floatingAddButton: () =>
      import('@/components/common/floating-add-button.vue')
  },
  computed: {
    tasksFiltered () {
      return this.tasks.filter((item) => {
        return this.TaskDisplayState(item)
      })
    }
  },
  methods: {
    GetTasks () {
      tasks
        .GetTasks(undefined, this.sortBy)
        .then((resp) => {
          this.pageLoading = false
          this.tasks = resp.data.list || []
        })
        .catch(() => {
          common.DisplayFailureToast(
            'Hmmm seems somethings gone wrong loading the tasks'
          )
        })
    },
    GetFlatmateName (id) {
      flatmates
        .GetFlatmate(id)
        .then((resp) => {
          return resp.data.spec.names
        })
        .catch((err) => {
          common.DisplayFailureToast(
            'Failed to fetch user account' +
              `<br/>${err.response.data.metadata.response}`
          )
          return id
        })
    },
    TaskDisplayState (list) {
      if (this.taskDisplayState === 1 && list.completed === false) {
        return this.ItemByNameInTask(list)
      } else if (this.taskDisplayState === 2 && list.completed === true) {
        return this.ItemByNameInTask(list)
      } else if (this.taskDisplayState === 0) {
        return this.ItemByNameInTask(list)
      }
    },
    ItemByNameInTask (item) {
      var vm = this
      return (
        item.name.toLowerCase().indexOf(vm.taskSearch.toLowerCase()) !== -1
      )
    },
    CheckDeviceIsMobile () {
      this.deviceIsMobile = common.DeviceIsMobile()
    }
  },
  watch: {
    sortBy () {
      this.taskIsLoading = true
      this.GetTasks()
    }
  },
  async beforeMount () {
    this.GetTasks()
  },
  beforeDestroy () {
    window.removeEventListener('resize', this.CheckDeviceIsMobile, true)
  },
  async created () {
    this.CheckDeviceIsMobile()
    window.addEventListener('resize', this.CheckDeviceIsMobile, true)
  }
}
</script>

<style scoped></style>
