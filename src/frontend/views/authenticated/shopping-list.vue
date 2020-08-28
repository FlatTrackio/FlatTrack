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
        <div v-if="notes !== '' || canUserAccountAdmin">
          <div class="content">
            <label class="label">Notes</label>
            <p :class="canUserAccountAdmin ? 'display-is-editable pointer-cursor-on-hover' : ''" class="subtitle is-4 notes-highlight" @click="EditShoppingListNotes">
              <i>
                {{ notes || "Add notes" }}
              </i>
            </p>
          </div>
          <br />
        </div>
        <b-button
          class="has-text-left"
          @click="goToRef('/apps/shopping-list/tags')"
          type="is-info"
          icon-left="tag-multiple"
          expanded>
          Manage tags
        </b-button>
        <br />
        <div>
          <b-tabs :position="deviceIsMobile ? 'is-centered' : ''" class="block is-marginless" v-model="listDisplayState">
            <b-tab-item icon="format-list-checks" label="All"></b-tab-item>
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
        <floatingAddButton path="/apps/shopping-list/new" v-if="displayFloatingAddButton"/>
        <br/>
        <div v-if="listsFiltered.length > 0">
          <shoppingListCardView :list="list" :authors="authors" :lists="lists" :index="index" v-for="(list, index) in listsFiltered" v-bind:key="list" :deviceIsMobile="deviceIsMobile" />
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
                  <p class="subtitle is-4" v-if="listSearch === '' && lists.length === 0 && !pageLoading">No lists added yet.</p>
                  <p class="subtitle is-4" v-else-if="listSearch === '' && listDisplayState === 1 && lists.length > 0 && !pageLoading">All lists have been completed.</p>
                  <p class="subtitle is-4" v-else-if="listSearch === '' && listDisplayState === 2 && lists.length > 0 && !pageLoading">No lists have been completed yet.</p>
                  <p class="subtitle is-4" v-else-if="listSearch !== '' && !pageLoading">No lists found.</p>
                  <p class="subtitle is-4" v-else-if="pageLoading">Loading lists...</p>
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
import cani from '@/frontend/requests/authenticated/can-i'
import { DialogProgrammatic as Dialog } from 'buefy'

export default {
  name: 'Shopping List',
  data () {
    return {
      displayFloatingAddButton: true,
      canUserAccountAdmin: false,
      notes: '',
      lists: [],
      authors: {},
      listDisplayState: 0,
      deviceIsMobile: false,
      listSearch: '',
      pageLoading: true,
      sortBy: 'recentlyUpdated'
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
    GetShoppingListNotes () {
      shoppinglist.GetShoppingListNotes().then(resp => {
        this.notes = resp.data.spec || ''
        this.pageLoading = false
      }).catch(() => {
        common.DisplayFailureToast('Hmmm seems somethings gone wrong loading the notes for shopping lists')
      })
    },
    EditShoppingListNotes () {
      if (this.canUserAccountAdmin !== true) {
        return
      }
      this.displayFloatingAddButton = false
      Dialog.prompt({
        title: 'Shopping list notes',
        message: `Enter notes that are useful for shopping in your flat.`,
        container: null,
        icon: 'text',
        hasIcon: true,
        inputAttrs: {
          placeholder: 'e.g. Our budget is $200/w. Please make sure to bring the supermarket card.',
          maxlength: 80,
          required: false,
          value: this.notes || undefined
        },
        trapFocus: true,
        onConfirm: (value) => {
          shoppinglist.PutShoppingListNotes(value).then(() => {
            this.pageLoading = true
            this.displayFloatingAddButton = true
            this.GetShoppingListNotes()
          }).catch(err => {
            common.DisplayFailureToast('Failed to update notes' + `<br/>${err.response.data.metadata.response}`)
            this.displayFloatingAddButton = true
          })
        },
        onCancel: () => {
          this.displayFloatingAddButton = true
        }
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
    cani.GetCanIgroup('admin').then(resp => {
      this.canUserAccountAdmin = resp.data.data
    })
    this.GetShoppingLists()
    this.GetShoppingListNotes()
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

<style scoped>

</style>
