<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><router-link to="/apps">Apps</router-link></li>
              <li class="is-active"><router-link to="/apps/shopping-list">Shopping list</router-link></li>
            </ul>
        </nav>
        <h1 class="title is-1">Shopping list</h1>
        <p class="subtitle is-3">Manage your weekly shop</p>
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
                    <p class="title is-3">Add a new list</p>
                  </div>
                </div>
              </div>
            </div>
          </section>
        </div>
        <br/>
        <div v-if="lists.length > 0">
          <shoppingListCardView :list="list" :authors="authors" v-for="list in lists" v-bind:key="list" />
          <br/>
          <p>{{ lists.length }} shopping list(s)</p>
        </div>
    </div>
  </div>
</template>

<script>
import common from '@/frontend/common/common'
import shoppinglist from '@/frontend/requests/authenticated/shoppinglist'
import flatmates from '@/frontend/requests/authenticated/flatmates'

export default {
  name: 'Shopping List',
  data () {
    return {
      lists: [],
      authors: {}
    }
  },
  components: {
    shoppingListCardView: () => import('@/frontend/components/authenticated/shopping-list-card-view.vue')
  },
  methods: {
    goToRef (ref) {
      this.$router.push({ path: ref })
    },
    GetShoppingLists () {
      shoppinglist.GetShoppingLists().then(resp => {
        this.lists = resp.data.list || []
      }).catch(() => {
        common.DisplayFailureToast('Hmmm seems somethings gone wrong loading the shopping lists')
      })
    },
    GetFlatmateName (id) {
      flatmates.GetFlatmate(id).then(resp => {
        return resp.data.spec.names
      }).catch(err => {
        common.DisplayFailureToast('Failed to fetch user account' + `<br/>${err.response.data.metadata.response}`)
        return id
      })
    }
  },
  async beforeMount () {
    this.GetShoppingLists()
  }
}
</script>

<style scoped>

</style>
