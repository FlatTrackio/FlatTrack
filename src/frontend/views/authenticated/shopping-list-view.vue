<template>
  <div>
    <div v-if="HeaderIsSticky && !editing" :class="ratherSmallerScreen ? 'ListBar ListBarTop' : 'ListBar'">
      <p class="subtitle is-5">
        <b>{{ name }}</b>
        ${{ currentPrice }}/${{ totalPrice }} ({{ Math.round(currentPrice / totalPrice * 100 * 100) / 100 || 0 }}%)
        <span @click="PatchShoppingListCompleted(id, !completed)" class="display-is-editable pointer-cursor-on-hover">
          <b-tag type="is-warning" v-if="!completed">Uncompleted</b-tag>
          <b-tag type="is-success" v-else-if="completed">Completed</b-tag>
        </span>
      </p>
    </div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
          <ul>
            <li><router-link to="/apps/shopping-list">Shopping list</router-link></li>
            <li v-if="hasInitialLoaded || name !== '' || typeof name !== 'undefined'" class="is-active"><router-link :to="'/apps/shopping-list/list/' + id">{{ name || 'Unnamed list' }}</router-link></li>
            <b-skeleton v-else size="is-small" width="35%" :animated="true"></b-skeleton>
          </ul>
        </nav>
        <div v-if="editingMeta">
          <b-field label="Name">
            <b-input
              type="text"
              icon="format-title"
              size="is-medium"
              ref="name"
              placeholder="Enter a title for this list"
              @keyup.enter.native="notesFromEmpty = false; editing = false; editingMeta = false; UpdateShoppingList(name, notes, completed)"
              @keyup.esc.native="editing = false; editingMeta = false"
              v-model="name"
              required>
            </b-input>
          </b-field>
          <br/>
        </div>
        <div v-else>
          <h1 v-if="hasInitialLoaded || name !== '' || typeof name !== 'undefined'" id="ListName" class="title is-1 is-marginless display-is-editable pointer-cursor-on-hover" @click="editing = true; editingMeta = true; FocusName()">{{ name }}</h1>
          <b-skeleton v-else size="is-medium" width="35%" :animated="true"></b-skeleton>
        </div>
        <div v-if="notes !== '' || notesFromEmpty || editingMeta">
          <div v-if="editingMeta">
            <b-field label="Notes">
              <b-input
                icon="text"
                size="is-medium"
                maxlength="100"
                type="text"
                ref="notes"
                placeholder="Enter extra information"
                @keyup.enter.native="notesFromEmpty = false; editing = false; editingMeta = false; UpdateShoppingList(name, notes)"
                @keyup.esc.native="editing = false; editingMeta = false"
                v-model="notes">
              </b-input>
            </b-field>
          </div>
          <div v-else>
            <br/>
            <div>
              <div class="content">
                <label class="label">Notes</label>
                <p class="display-is-editable subtitle is-4 pointer-cursor-on-hover notes-highlight" @click="editing = true; editingMeta = true; FocusNotes()">
                  <i>
                    {{ notes }}
                  </i>
                </p>
              </div>
            </div>
          </div>
        </div>
        <b-button type="is-text" @click="() => { notesFromEmpty = true; editing = true; editingMeta = true; FocusNotes() }" v-if="!editingMeta && notes.length == 0" class="remove-shadow">Add notes</b-button>
        <div v-if="editingMeta">
          <b-button type="is-info" @click="() => { notesFromEmpty = false; editing = false; editingMeta = false; UpdateShoppingList(name, notes, completed) }">Done</b-button>
          <br/>
        </div>
        <br/>
        <b-tag type="is-success" v-if="completed">Completed</b-tag>
        <b-tag type="is-warning" v-else-if="!completed">Uncompleted</b-tag>
        <br/>
        <b-tabs :position="deviceIsMobile ? 'is-centered' : ''" class="block is-marginless" v-model="itemDisplayState">
          <b-tab-item icon="format-list-checks" label="All"></b-tab-item>
          <b-tab-item icon="playlist-remove" label="Unobtained"></b-tab-item>
          <b-tab-item icon="playlist-check" label="Obtained"></b-tab-item>
        </b-tabs>
        <label class="label">Search for items</label>
        <b-field>
          <b-input
            icon="magnify"
            size="is-medium"
            placeholder="Item name"
            type="search"
            v-model="itemSearch"
            ref="search"
            v-on:keyup.ctrl.66="FocusSearchBox"
            expanded>
          </b-input>
          <p class="control">
            <b-select
              placeholder="Sort by"
              icon="sort"
              v-model="sortBy"
              size="is-medium"
              expanded>
              <option value="tags">Tags</option>
              <option value="highestPrice">Highest Price</option>
              <option value="lowestPrice">Lowest Price</option>
              <option value="highestQuantity">Highest Quantity</option>
              <option value="lowestQuantity">Lowest Quantity</option>
              <option value="recentlyAdded">Recently Added</option>
              <option value="lastAdded">Last Added</option>
              <option value="recentlyUpdated">Recently Updated</option>
              <option value="lastUpdated">Last Updated</option>
              <option value="alphabeticalDescending">A-z</option>
              <option value="alphabeticalAscending">z-A</option>
            </b-select>
          </p>
        </b-field>
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
                  <div class="media-right">
                    <b-icon icon="chevron-right" size="is-medium" type="is-midgray"></b-icon>
                  </div>
                </div>
              </div>
            </div>
          </section>
        </div>
        <br/>
        <div v-if="listItemsFromTags.length > 0">
          <b-loading :is-full-page="false" :active.sync="listIsLoading" :can-cancel="false"></b-loading>
          <div v-if="sortBy === 'tags'">
            <section v-for="itemTag in listItemsFromTags" v-bind:key="itemTag">
              <div v-if="editingTag === itemTag.tag">
                <label class="label">Tag name</label>
                <b-field>
                  <b-button
                    type="is-danger"
                    size="is-medium"
                    @click="editingTag = ''; editing = false">
                    X
                  </b-button>
                  <b-input
                    type="text"
                    icon="format-title"
                    size="is-medium"
                    maxlength="30"
                    placeholder="Enter a tag name"
                    expanded
                    @keyup.enter.native="editingTag = ''; UpdateShoppingListItemTag(itemTag.tag, TagTmp); itemTag.tag = TagTmp; TagTmp = ''; editing = false"
                    @keyup.esc.native="editingTag = ''; editing = false"
                    v-model="TagTmp"
                    required>
                  </b-input>
                  <p class="control">
                    <b-button
                      type="is-primary"
                      size="is-medium"
                      icon-left="check"
                      @click="editingTag = ''; UpdateShoppingListItemTag(itemTag.tag, TagTmp); itemTag.tag = TagTmp; TagTmp = ''; editing = false">
                    </b-button>
                  </p>
                </b-field>
                <br/>
              </div>
              <div v-else>
                <div class="field">
                  <p class="control">
                    <b-button
                      type="is-text"
                      size="medium"
                      class="title is-5 is-marginless display-is-editable pointer-cursor-on-hover is-paddingless remove-shadow"
                      @click="TagTmp = itemTag.tag; editingTag = itemTag.tag; editing = true">
                      {{ itemTag.tag }}
                    </b-button>
                    <b-button
                      style="float: right; display: block;"
                      icon-left="plus-box"
                      size="medium"
                      tag="router-link"
                      :to="{ name: 'New shopping list item', query: { 'tag': itemTag.tag }}">
                    </b-button>
                  </p>
                </div>
              </div>
              <transition-group
                name="staggered-fade"
                tag="div"
                v-bind:css="false"
                v-on:enter="ItemAppear"
                v-on:leave="ItemDisappear">
                <div v-for="(item, index) in itemTag.items" v-bind:key="item">
                  <itemCard :list="list" :item="item" :index="index" :listId="id" :deviceIsMobile="deviceIsMobile"/>
                </div>
                <br/>
              </transition-group>
              <section>
                <br/>
                <p>
                  {{ itemTag.items.length || 0 }} item(s)
                  <span v-if="itemTag.price !== 0 && typeof itemTag.price !== 'undefined'">
                    - ${{ itemTag.price.toFixed(2) }}
                  </span>
                </p>
              </section>
              <br/>
            </section>
          </div>
          <div v-else-if="sortBy !== 'tag'">
            <div v-for="(item, index) in listItemsFromPlainList" v-bind:key="item">
              <itemCard :list="list" :item="item" :index="index" :listId="id" :displayTag="true" :deviceIsMobile="deviceIsMobile"/>
            </div>
            <section>
              <br/>
              <p>
                {{ listItemsFromPlainList.length || 0 }} item(s)
              </p>
            </section>
            <br/>
          </div>
        </div>
        <div v-else>
          <div class="card">
            <div class="card-content card-content-list">
              <div class="media">
                <div class="media-left">
                  <b-icon icon="cart-remove" size="is-medium" type="is-midgray"></b-icon>
                </div>
                <div class="media-content">
                  <p class="subtitle is-4" v-if="itemSearch === '' && list.length === 0">No items added yet.</p>
                  <p class="subtitle is-4" v-else-if="itemSearch === '' && itemDisplayState === 1 && list.length > 0">All items have been obtained.</p>
                  <p class="subtitle is-4" v-else-if="itemSearch === '' && itemDisplayState === 2 && list.length > 0">No items have been obtained yet.</p>
                  <p class="subtitle is-4" v-else-if="itemSearch !== ''">No items found.</p>
                </div>
              </div>
            </div>
          </div>
          <br/>
        </div>
        <floatingAddButton :path="'/apps/shopping-list/list/' + id + '/new'"/>
        <p class="subtitle is-4">
          <b>Total items</b>: {{ obtainedCount }}/{{ totalItems }}
          <br/>
          <b>Total price</b>: ${{ currentPrice }}/${{ totalPrice }} ({{ Math.round(currentPrice / totalPrice * 100 * 100) / 100 || 0 }}%)
        </p>
        <b-field>
          <b-button
            :icon-left="completed === false ? 'checkbox-blank-outline' : 'check-box-outline'"
            :type="completed === true ? 'is-success' : 'is-warning'"
            size="is-medium"
            expanded
            @click="PatchShoppingListCompleted(id, !completed)">
            {{ completed === false ? 'Uncompleted' : 'Completed' }}
          </b-button>
          <p class="control">
            <b-button
              icon-left="delete"
              type="is-danger"
              size="is-medium"
              :loading="deleteLoading"
              @click="DeleteShoppingList(id)">
            </b-button>
          </p>
        </b-field>
        <p class="subtitle is-6">
          Created
          <span v-if="hasInitialLoaded || typeof authorName !== 'undefined'">{{ TimestampToCalendar(creationTimestamp) }}, by</span>
          <b-skeleton v-else size="is-small" width="35%" :animated="true"></b-skeleton>
          <router-link v-if="hasInitialLoaded || typeof authorName !== 'undefined'" tag="a" :to="{ name: 'My Flatmates', query: { 'id': author }}"> {{ authorNames }} </router-link>
            <b-skeleton v-else size="is-small" width="35%" :animated="true"></b-skeleton>
            <span v-if="creationTimestamp !== modificationTimestamp">
              <br/>
              Last updated
              <span v-if="hasInitialLoaded || typeof authorName !== 'undefined'">{{ TimestampToCalendar(modificationTimestamp) }}, by</span>
              <b-skeleton v-else size="is-small" width="35%" :animated="true"></b-skeleton>
              <router-link v-if="hasInitialLoaded || typeof authorName !== 'undefined'" tag="a" :to="{ name: 'My Flatmates', query: { 'id': author }}"> {{ authorLastNames }} </router-link>
            <b-skeleton v-else size="is-small" width="35%" :animated="true"></b-skeleton>
            </span>
        </p>
        <br/>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/frontend/common/common'
import shoppinglistCommon from '@/frontend/common/shoppinglist'
import shoppinglist from '@/frontend/requests/authenticated/shoppinglist'
import flatmates from '@/frontend/requests/authenticated/flatmates'
import { DialogProgrammatic as Dialog } from 'buefy'

export default {
  name: 'Shopping List',
  data () {
    return {
      intervalLoop: null,
      editing: false,
      editingMeta: false,
      notesFromEmpty: false,
      itemSearch: '',
      authorNames: '',
      authorLastNames: '',
      totalItems: 0,
      loopCreated: new Date(),
      sortBy: shoppinglistCommon.GetShoppingListSortBy() || 'tags',
      itemDisplayState: null,
      deviceIsMobile: false,
      HeaderIsSticky: false,
      TagTmp: '',
      editingTag: '',
      listIsLoading: true,
      hasInitialLoaded: false,
      deleteLoading: false,
      ratherSmallerScreen: false,
      id: this.$route.params.id,
      name: 'Unnamed list',
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
    itemCard: () => import('@/frontend/components/authenticated/shopping-list-item-card-view.vue'),
    floatingAddButton: () => import('@/frontend/components/common/floating-add-button.vue')
  },
  computed: {
    listItemsFromTags () {
      return this.RestructureShoppingListToTags(this.list.filter((item) => {
        return this.ItemDisplayState(item)
      }))
    },
    listItemsFromPlainList () {
      return this.list.filter((item) => {
        return this.ItemDisplayState(item)
      })
    },
    obtainedCount () {
      var obtained = 0
      var list = this.RestructureShoppingListToTags(this.list)
      for (var itemTag in list) {
        for (var item in list[itemTag].items) {
          if (list[itemTag].items[item].obtained === true) {
            obtained += 1
          }
        }
      }
      return obtained
    },
    currentPrice () {
      var currentPrice = 0
      var list = this.RestructureShoppingListToTags(this.list)
      for (var itemTag in list) {
        for (var item in list[itemTag].items) {
          var itemInList = list[itemTag].items[item]
          itemInList.price = typeof itemInList.price === 'undefined' ? 0 : itemInList.price
          if (itemInList.obtained === true) {
            currentPrice += itemInList.price * itemInList.quantity
          }
        }
      }
      currentPrice = Math.round(currentPrice * 100) / 100
      return currentPrice
    },
    totalPrice () {
      var totalPrice = 0
      var list = this.RestructureShoppingListToTags(this.list)
      for (var itemTag in list) {
        for (var item in list[itemTag].items) {
          var itemInList = list[itemTag].items[item]
          itemInList.price = typeof itemInList.price === 'undefined' ? 0 : itemInList.price
          totalPrice += itemInList.price * itemInList.quantity
        }
      }
      totalPrice = Math.round(totalPrice * 100) / 100
      return totalPrice
    }
  },
  methods: {
    ItemByNameInList (item) {
      var vm = this
      return item.name.toLowerCase().indexOf(vm.itemSearch.toLowerCase()) !== -1
    },
    ItemDisplayState (item) {
      var vm = this
      if (this.itemDisplayState === 1 && item.obtained === false) {
        return this.ItemByNameInList(item)
      } else if (this.itemDisplayState === 2 && item.obtained === true) {
        return this.ItemByNameInList(item)
      } else if (this.itemDisplayState === 0) {
        return this.ItemByNameInList(item)
      }
    },
    goToRef (ref) {
      this.$router.push({ path: ref })
    },
    FocusSearchBox () {
      this.$refs.search.$el.focus()
    },
    RestructureShoppingListToTags (list) {
      return shoppinglistCommon.RestructureShoppingListToTags(list)
    },
    GetShoppingList () {
      if (this.editing === true) {
        return
      }
      var id = this.id
      shoppinglist.GetShoppingList(id).then(resp => {
        this.name = resp.data.spec.name
        this.notes = resp.data.spec.notes || ''
        this.author = resp.data.spec.author
        this.authorLast = resp.data.spec.authorLast
        this.completed = resp.data.spec.completed
        this.creationTimestamp = resp.data.spec.creationTimestamp
        this.modificationTimestamp = resp.data.spec.modificationTimestamp
        return flatmates.GetFlatmate(this.author)
      }).then(resp => {
        this.authorNames = resp.data.spec.names
        return flatmates.GetFlatmate(this.authorLast)
      }).then(resp => {
        this.authorLastNames = resp.data.spec.names
      }).catch(err => {
        if (err.response.status === 404) {
          common.DisplayFailureToast('Error list not found' + '<br/>' + err.response.data.metadata.response)
          this.$router.push({ name: 'Shopping list' })
          return
        }
        common.DisplayFailureToast('Error loading the shopping list' + '<br/>' + err.response.data.metadata.response)
      })
    },
    UpdateShoppingList (name, notes, completed) {
      var id = this.id
      shoppinglist.UpdateShoppingList(id, name, notes, completed).catch(err => {
        common.DisplayFailureToast('Failed to update shopping list' + '<br/>' + err.response.data.metadata.response)
      })
    },
    PatchShoppingListCompleted (id, completed) {
      shoppinglist.PatchShoppingListCompleted(id, completed).then(resp => {
        this.completed = resp.data.spec.completed
      }).catch(err => {
        common.DisplayFailureToast('Failed to set list as completed' + '<br/>' + err.response.data.metadata.response)
      })
    },
    DeleteShoppingList (id) {
      Dialog.confirm({
        title: 'Delete shopping list',
        message: 'Are you sure that you wish to delete this shopping list?' + '<br/>' + 'This action cannot be undone.',
        confirmText: 'Delete shopping list',
        type: 'is-danger',
        hasIcon: true,
        onConfirm: () => {
          this.deleteLoading = true
          window.clearInterval(this.intervalLoop)
          shoppinglist.DeleteShoppingList(id).then(resp => {
            common.DisplaySuccessToast('Deleted the shopping list')
            shoppinglistCommon.DeleteShoppingListFromCache(id)
            setTimeout(() => {
              this.$router.push({ name: 'Shopping list' })
            }, 1 * 1000)
          }).catch(err => {
            this.deleteLoading = false
            common.DisplayFailureToast('Failed to delete the shopping list' + '<br/>' + err.response.data.metadata.response)
          })
        }
      })
    },
    GetShoppingListItems () {
      shoppinglist.GetShoppingListItems(this.id, this.sortBy).then(resp => {
        var responseList = resp.data.list
        this.totalItems = responseList === null ? 0 : responseList.length
        if (this.list === null) {
          this.list = []
        }

        if (responseList !== this.list) {
          this.list = responseList || []
          shoppinglistCommon.WriteShoppingListToCache(this.id, this.list)
          this.listIsLoading = false
          this.hasInitialLoaded = true
        }
      })
    },
    UpdateShoppingListItemTag (tagName, tagNameNew) {
      shoppinglist.UpdateShoppingListItemTag(this.id, tagName, tagNameNew).catch(err => {
        common.DisplayFailureToast('Failed to update the shopping list tag' + '<br/>' + err.response.data.metadata.response)
      })
    },
    ItemAppear (el, done) {
      var delay = el.dataset.index * 150
      setTimeout(function () {
        Velocity(
          el,
          { opacity: 1, height: '1.6em' },
          { complete: done }
        )
      }, delay)
    },
    ItemDisappear (el, done) {
      var delay = el.dataset.index * 150
      setTimeout(function () {
        Velocity(
          el,
          { opacity: 0, height: 0 },
          { complete: done }
        )
      }, delay)
    },
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    },
    LoopStart () {
      if (shoppinglistCommon.GetShoppingListAutoRefresh() === 'false') {
        return
      }
      this.intervalLoop = window.setInterval(() => {
        if (this.editing === true) {
          return
        }
        this.GetShoppingList()
        this.GetShoppingListItems()

        var now = new Date()
        var timePassed = (now.getTime() / 1000) - (this.loopCreated.getTime() / 1000)
        if (timePassed >= 3600 / 4) {
          window.clearInterval(this.intervalLoop)
        }
      }, 3 * 1000)
    },
    LoopStop () {
      window.clearInterval(this.intervalLoop)
    },
    CheckDeviceIsMobile () {
      this.deviceIsMobile = common.DeviceIsMobile()
    },
    ManageStickyHeader () {
      this.HeaderIsSticky = window.pageYOffset > document.getElementById('ListName').offsetTop + 30
    },
    ResetLoopTime () {
      this.loopCreated = new Date()
    },
    FocusName () {
      this.$refs.name.$el.focus()
    },
    FocusNotes () {
      this.$refs.notes.$el.focus()
    }
  },
  watch: {
    sortBy () {
      shoppinglistCommon.WriteShoppingListSortBy(this.sortBy)
      this.listIsLoading = true
      this.ResetLoopTime()
      this.LoopStop()
      this.LoopStart()
    },
    itemDisplayState () {
      shoppinglistCommon.WriteShoppingListObtainedFilter(this.id, this.itemDisplayState)
    },
    itemSearch () {
      shoppinglistCommon.WriteShoppingListSearch(this.id, this.itemSearch)
    }
  },
  async beforeMount () {
    this.list = shoppinglistCommon.GetShoppingListFromCache(this.id) || []
    this.itemSearch = shoppinglistCommon.GetShoppingListSearch(this.id) || ''
    this.GetShoppingList()
    this.GetShoppingListItems()
    if (window.innerWidth <= 330) {
      this.ratherSmallerScreen = true
    }
  },
  async created () {
    this.itemDisplayState = shoppinglistCommon.GetShoppingListObtainedFilter(this.id) || 0
    this.CheckDeviceIsMobile()
    window.addEventListener('resize', this.CheckDeviceIsMobile, true)
    window.addEventListener('scroll', this.ManageStickyHeader, true)
    this.LoopStart()
    window.addEventListener('focus', this.ResetLoopTime, true)
  },
  beforeDestroy () {
    this.LoopStop()
    window.removeEventListener('resize', this.CheckDeviceIsMobile, true)
    window.removeEventListener('scroll', this.ManageStickyHeader, true)
    window.removeEventListener('focus', this.ResetLoopTime, true)
  }
}
</script>

<style scoped>
.display-is-editable:hover {
    text-decoration: underline dotted;
    -webkit-transition: width 0.5s ease-in;
}
.card-content-list {
    background-color: transparent;
    padding-left: 1.5em;
    padding-top: 0.6em;
    padding-bottom: 0.6em;
    padding-right: 1.5em;
}

.obtained {
    color: #adadad;
    text-decoration: line-through;
}

.ListBar {
    position: fixed;
    height: auto;
    width: 100%;
    z-index: 20;
    padding: 10px;
    box-shadow: black 0px -45px 71px;
    display: block;
    background-color: hsla(0,0%,100%,.73);
    backdrop-filter: blur(5px);
}

.ListBarTop {
    top: 0;
}
</style>
