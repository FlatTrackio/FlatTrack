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
          <b-tabs :position="deviceIsMobile ? 'is-centered' : ''" class="block is-marginless" v-model="listDisplayState">
            <b-tab-item icon="" label="All"></b-tab-item>
            <b-tab-item icon="playlist-remove" label="Uncompleted"></b-tab-item>
            <b-tab-item icon="playlist-check" label="Completed"></b-tab-item>
          </b-tabs>
          <label class="label">Search for lists</label>
          <b-field>
            <b-input
              icon="magnify"
              size="is-medium"
              placeholder="Enter a list name"
              type="search"
              expanded
              v-model="listSearch"
              ref="search">
            </b-input>
            <p class="control">
              <b-select
                placeholder="Sort by"
                icon="sort"
                v-model="sortBy"
                size="is-medium"
                expanded>
                <option value="recentlyAdded">Recently Added</option>
                <option value="lastAdded">Last Added</option>
                <option value="recentlyUpdated">Recently Updated</option>
                <option value="lastUpdated">Last Updated</option>
                <option value="alphabeticalDescending">A-z</option>
                <option value="alphabeticalAscending">z-A</option>
            </b-select>
          </p>
          </b-field>
          <b-loading :is-full-page="false" :active.sync="pageLoading" :can-cancel="false"></b-loading>
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
        <div v-else>
          <div class="card">
            <div class="card-content card-content-list">
              <div class="media">
                <div class="media-left">
                  <b-icon icon="cart-remove" size="is-medium" type="is-midgray"></b-icon>
                </div>
                <div class="media-content">
                  <p class="subtitle is-4" v-if="listSearch === '' && lists.length === 0">No lists added yet.</p>
                  <p class="subtitle is-4" v-else-if="listSearch === '' && listDisplayState === 1 && lists.length > 0">All lists have been completed.</p>
                  <p class="subtitle is-4" v-else-if="listSearch === '' && listDisplayState === 2 && lists.length > 0">No lists have been completed yet.</p>
                  <p class="subtitle is-4" v-else-if="listSearch !== ''">No lists found.</p>
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
      deviceIsMobile: false,
      listSearch: '',
      pageLoading: true,
      sortBy: 'lastUpdated'
    }
  },
  components: {
    shoppingListCardView: () => import('@/frontend/components/authenticated/shopping-list-card-view.vue'),
    floatingAddButton: () => import('@/frontend/components/common/floating-add-button.vue')
  },
  computed: {
    listsFiltered () {
      return this.lists.filter((item) => {
        return this.ListDisplayState(item)
      })
    }
  },
  methods: {
    goToRef (ref) {
      this.$router.push({ path: ref })
    },
    GetShoppingLists () {
      shoppinglist.GetShoppingLists(undefined, this.sortBy).then(resp => {
        this.pageLoading = false
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
    ListDisplayState (list) {
      var vm = this
      if (this.listDisplayState === 1 && list.completed === false) {
        return this.ItemByNameInList(list)
      } else if (this.listDisplayState === 2 && list.completed === true) {
        return this.ItemByNameInList(list)
      } else if (this.listDisplayState === 0) {
        return this.ItemByNameInList(list)
      }
    },
    ItemByNameInList (item) {
      var vm = this
      return item.name.toLowerCase().indexOf(vm.listSearch.toLowerCase()) !== -1
    },
    CheckDeviceIsMobile () {
      this.deviceIsMobile = common.DeviceIsMobile()
    }
  },
  watch: {
    sortBy () {
      this.listIsLoading = true
      this.GetShoppingLists()
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
