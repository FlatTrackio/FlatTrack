<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
          <ul>
            <li><router-link to="/apps">Apps</router-link></li>
            <li><router-link to="/apps/shopping-list">Apps</router-link></li>
            <li class="is-active"><router-link :to="'/apps/shopping-list/list/' + id">Shopping list</router-link></li>
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
        <div v-if="notes">
          <div v-if="editing">
            <b-field>
              <b-input maxlength="100" type="textarea" v-model="notes" class="subtitle is-3"></b-input>
            </b-field>
          </div>
          <div v-else>
            <div class="notification">
              <div class="content">
                <p class="display-is-editable" @click="editing = !editing">
                  {{ notes }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <b-button type="is-info" @click="editing = !editing" v-if="editing">Done</b-button>
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
  }
}
</script>

<style scoped>
.display-is-editable:hover {
    text-decoration: underline dotted;
}
</style>
