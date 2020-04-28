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
          <b-tabs :position="deviceIsMobile ? 'is-centered' : ''" class="block" v-model="listDisplayState">
            <b-tab-item icon="" label="All"></b-tab-item>
            <b-tab-item icon="playlist-remove" label="Uncompleted"></b-tab-item>
            <b-tab-item icon="playlist-check" label="Completed"></b-tab-item>
          </b-tabs>
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
                  <div class="media-right">
                    <b-icon icon="chevron-right" size="is-medium" type="is-midgray"></b-icon>
                  </div>
                </div>
              </div>
            </div>
          </section>
        </div>
        <floatingAddButton path="/apps/shopping-list/new"/>
        <br/>
        <div v-if="listsFiltered.length > 0">
          <shoppingListCardView :list="list" :authors="authors" v-for="list in listsFiltered" v-bind:key="list" />
          <br/>
          <p>{{ listsFiltered.length }} shopping list(s)</p>
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
      authors: {},
      listDisplayState: 0,
      deviceIsMobile: false
    }
  },
  components: {
    shoppingListCardView: () => import('@/frontend/components/authenticated/shopping-list-card-view.vue'),
    floatingAddButton: () => import('@/frontend/components/common/floating-add-button.vue')
  },
  computed: {
    listsFiltered () {
      return this.lists.filter((item) => {
        if (this.listDisplayState === 1 && item.completed === false) {
          return item
        } else if (this.listDisplayState === 2 && item.completed === true) {
          return item
        } else if (this.listDisplayState === 0) {
          return item
        }
      })
    }
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
    },
    CheckDeviceIsMobile () {
      this.deviceIsMobile = common.DeviceIsMobile()
    }
  },
  async beforeMount () {
    this.GetShoppingLists()
  },
  async created () {
    this.CheckDeviceIsMobile()
    window.addEventListener('resize', this.CheckDeviceIsMobile.bind(this))
  }
}
</script>

<style scoped>

</style>
