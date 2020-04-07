<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
          <ul>
            <li><router-link to="/apps">Apps</router-link></li>
            <li><router-link to="/apps/shopping-list">Shopping list</router-link></li>
            <li class="is-active"><router-link :to="'/apps/shopping-list/list/' + id">{{ name }}</router-link></li>
          </ul>
        </nav>
        <nav class="level">
          <div class="level-left">
            <div v-if="editing">
              <div class="field">
                <div class="control is-clearfix">
                  <input type="text" autocomplete="on" class="input title is-1" v-model="name"/>
                </div>
              </div>
            </div>
            <div v-else>
              <h1 class="title is-1 display-is-editable" @click="editing = !editing">{{ name }}</h1>
            </div>
          </div>
        </nav>
        <div v-if="!editing">
          <p>
          <span v-if="creationTimestamp == modificationTimestamp">
            Created
          </span>
          <span v-else>
            Updated
          </span>
          {{ TimestampToCalendar(creationTimestamp) }}, by {{ author }}
          </p>
          <br/>
        </div>
        <div v-if="notes || notesFromEmpty || editing">
          <div v-if="editing">
            <b-field>
              <b-input maxlength="100" :type="notes.length > 30 ? 'textarea' : 'text'" v-model="notes" class="subtitle is-3"></b-input>
            </b-field>
          </div>
          <div v-else>
            <div class="notification">
              <div class="content">
                <p class="display-is-editable subtitle is-4" @click="editing = true">
                  {{ notes }}
                </p>
              </div>
            </div>
          </div>
        </div>
        <b-button type="is-info" @click="() => { notesFromEmpty = true; editing = true }" v-if="!editing && notes.length == 0">Add notes</b-button>
        <b-button type="is-info" @click="editing = false" v-if="editing">Done</b-button>
        <br/>
        <br/>
        <div>
          <section>
            <div class="card pointer-cursor-on-hover" @click="goToRef('/apps/shopping-list/list/' + id + '/new')">
              <div class="card-content">
                <div class="media">
                  <div class="media-left">
                    <b-icon
                      icon="plus-box"
                      size="is-medium">
                    </b-icon>
                  </div>
                  <div class="media-content">
                    <p class="subtitle is-4">Add a new item</p>
                  </div>
                </div>
              </div>
            </div>
          </section>
        </div>
        <br/>
        <shoppinglistTable :listId="id" :items="list"/>
        <br/>
        <p>
          <b>Total items</b>: {{ list.length }}
        </p>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/frontend/common/common'
import shoppinglist from '@/frontend/requests/authenticated/shoppinglist'

export default {
  name: 'Shopping List',
  data () {
    return {
      editing: false,
      notesFromEmpty: false,
      id: this.$route.params.id,
      name: '',
      notes: '',
      author: '',
      authorLast: '',
      completed: false,
      creationTimestamp: 0,
      modificationTimestamp: 0,
      list: []
    }
  },
  components: {
    shoppinglistTable: () => import('@/frontend/components/authenticated/shopping-list.vue')
  },
  methods: {
    goToRef (ref) {
      this.$router.push({ path: ref })
    },
    GetShoppingList () {
      var id = this.id
      shoppinglist.GetShoppingList(id).then(resp => {
        console.log(resp)
        this.name = resp.data.spec.name
        this.notes = resp.data.spec.notes
        this.author = resp.data.spec.author
        this.authorLast = resp.data.spec.authorLast
        this.completed = resp.data.spec.completed
        this.creationTimestamp = resp.data.spec.creationTimestamp
        this.modificationTimestamp = resp.data.spec.modificationTimestamp
      }).catch(() => {
        common.DisplayFailureToast('Hmmm seems somethings gone wrong loading the shopping list')
      })
    },
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    }
  },
  async created () {
    this.GetShoppingList()
    shoppinglist.GetShoppingListItems(this.id).then(resp => {
      this.list = resp.data.list
      if (this.list === null) {
        this.list = []
      }
      console.log(this.list)
    })
  }
}
</script>

<style scoped>
.display-is-editable:hover {
    text-decoration: underline dotted;
}
</style>
