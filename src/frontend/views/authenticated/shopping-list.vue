<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><router-link to="/">Home</router-link></li>
              <li><router-link to="/apps">Apps</router-link></li>
              <li class="is-active"><router-link to="/apps/shopping-list">Shopping list</router-link></li>
            </ul>
        </nav>
        <h1 class="title is-1">Shopping list</h1>
        <h2 class="subtitle">Manage your weekly shop</h2>
        <b-tabs position="is-centered" class="block">
          <b-tab-item label="In Progress"></b-tab-item>
          <b-tab-item label="Completed" :disabled="lists.length === 0"></b-tab-item>
        </b-tabs>
        <div>
          <section>
            <div class="card pointer-cursor-on-hover" @click="goToRef('/apps/shopping-list/new')">
              <div class="card-content">
                <div class="media">
                  <div class="media-left">
                    <b-icon
                      icon="cart-plus"
                      size="is-medium">
                    </b-icon>
                  </div>
                  <div class="media-content">
                    <p class="title is-4">Add a new list</p>
                  </div>
                </div>
              </div>
            </div>
          </section>
        </div>
        <br/>
        <div v-if="lists.length > 0">
          <div v-for="list in lists" v-bind:key="list">
            <section>
              <div class="card pointer-cursor-on-hover" @click="goToRef('/apps/shopping-list/list/' + list.id)">
                <div class="card-content">
                  <div class="media">
                    <div class="media-left">
                      <b-icon
                        icon="cart-outline"
                        size="is-medium">
                      </b-icon>
                    </div>
                    <div class="media-content">
                      <p class="title is-4">{{ list.name }}</p>
                      <p class="subtitle is-6">
                        <span v-if="list.creationTimestamp == list.modificationTimestamp">
                          Created
                        </span>
                        <span v-else>
                          Updated
                        </span>
                        {{ TimestampToCalendar(list.creationTimestamp) }}, by {{ list.author }}
                      </p>
                    </div>
                  </div>
                  <div class="content" v-if="list.notes">
                    {{ list.notes }}
                    <br/>
                    <br/>
                    <div v-if="typeof list.count !== 'undefined' && list.count > 0">
                      {{ list.count }} item(s)
                    </div>
                    <div v-else>
                      0 items
                    </div>
                  </div>
                </div>
              </div>
            </section>
          </div>
          <br/>
          <p>{{ lists.length }} shopping list(s)</p>
        </div>
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
      lists: []
    }
  },
  methods: {
    goToRef (ref) {
      this.$router.push({ path: ref })
    },
    GetShoppingLists () {
      shoppinglist.GetShoppingLists().then(resp => {
        this.lists = resp.data.list
        console.log(this.lists)
        if (this.lists === null) {
          this.lists = []
        }
      }).catch(() => {
        common.DisplayFailureToast('Hmmm seems somethings gone wrong loading the shopping lists')
      })
    },
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    }
  },
  async created () {
    this.GetShoppingLists()
  }
}
</script>

<style scoped>

</style>
