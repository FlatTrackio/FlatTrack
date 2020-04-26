<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
          <ul>
            <li><router-link to="/apps/shopping-list">Shopping list</router-link></li>
            <li class="is-active"><router-link :to="'/apps/shopping-list/list/' + id">{{ name || 'Unnamed list' }}</router-link></li>
          </ul>
        </nav>
        <div v-if="editing">
          <b-field label="Name">
            <b-input
              type="text"
              icon="format-title"
              size="is-medium"
              @keyup.enter.native="notesFromEmpty = false; editing = false; UpdateShoppingList(name, notes, completed)"
              v-model="name"
              required>
            </b-input>
          </b-field>
          <br/>
        </div>
        <div v-else>
          <h1 class="title is-1 is-marginless display-is-editable pointer-cursor-on-hover" @click="editing = !editing">{{ name || 'Unnamed list' }}</h1>
          <br/>
          <b-tag type="is-info" v-if="completed">Completed</b-tag>
          <b-tag type="is-warning" v-if="!completed">Uncompleted</b-tag>
          <br/>
          <br/>
          <p class="subtitle is-6">
            Created {{ TimestampToCalendar(creationTimestamp) }}, by <router-link tag="a" :to="'/apps/flatmates?id=' + author"> {{ authorNames }} </router-link>
            <span v-if="creationTimestamp !== modificationTimestamp">
              <br/>
              Last updated {{ TimestampToCalendar(modificationTimestamp) }}, by <router-link tag="a" :to="'/apps/flatmates?id=' + authorLast"> {{ authorLastNames }} </router-link>
            </span>
          </p>
        </div>
        <div v-if="notes != '' || notesFromEmpty || editing">
          <div v-if="editing">
            <b-field label="Notes">
              <b-input
                icon="text"
                size="is-medium"
                maxlength="100"
                type="text"
                @keyup.enter.native="notesFromEmpty = false; editing = false; UpdateShoppingList(name, notes)"
                v-model="notes">
              </b-input>
            </b-field>
          </div>
          <div v-else>
            <br/>
            <div class="notification">
              <div class="content">
                <label class="label">Notes</label>
                <p class="display-is-editable subtitle is-4 pointer-cursor-on-hover" @click="editing = true">
                  {{ notes }}
                </p>
              </div>
            </div>
          </div>
        </div>
        <b-button type="is-text" @click="() => { notesFromEmpty = true; editing = true }" v-if="!editing && notes.length == 0">Add notes</b-button>
        <div v-if="editing">
          <b-button type="is-info" @click="() => { notesFromEmpty = false; editing = false; UpdateShoppingList(name, notes, completed) }">Done</b-button>
          <br/>
        </div>
        <br/>
        <label class="label">Search for items</label>
        <b-field>
          <b-input icon="magnify" size="is-medium" placeholder="Item name" type="search" v-model="itemSearch" ref="search" v-on:keyup.ctrl.66="FocusSearchBox"></b-input>
          <p class="control">
            <b-select placeholder="Sort by" icon="sort" v-model="sortBy" size="is-medium" expanded>
              <option value="tags">Tags</option>
              <option value="price">Price</option>
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
        <div v-if="listItems.length > 0">
          <div v-if="sortBy === 'tags'">
            <section v-for="itemTag in listItems" v-bind:key="itemTag">
              <p class="title is-5">
                {{ itemTag.tag }}
                <span v-if="itemTag.price !== 0 && typeof itemTag.price !== 'undefined'">
                  (${{ itemTag.price }})
                </span>
              </p>
              <transition-group
                name="staggered-fade"
                tag="div"
                v-bind:css="false"
                v-on:enter="ItemAppear"
                v-on:leave="ItemDisappear">
                <div v-for="(item, index) in itemTag.items" v-bind:key="item">
                  <itemCard :list="list" :item="item" :index="index" :listId="id"/>
                </div>
                <br/>
              </transition-group>
              <section>
                <br/>
                <p>
                  {{ itemTag.items.length || 0 }} item(s)
                </p>
              </section>
              <br/>
            </section>
          </div>
          <div v-if="sortBy === 'price'">
            <div v-for="(item, index) in list" v-bind:key="item">
              <itemCard :list="list" :item="item" :index="index" :listId="id" :displayTag="true"/>
            </div>
          </div>
          <br/>
        </div>
        <div v-else>
          <div class="card">
            <div class="card-content card-content-list">
              <div class="media">
                <div class="media-left" @click="PatchItemObtained(item.id, !item.obtained)">
                  <b-icon icon="cart-remove" size="is-medium" type="is-midgray"></b-icon>
                </div>
                <div class="media-content">
                  <p class="subtitle is-4" v-if="itemSearch === ''">No items added yet.</p>
                  <p class="subtitle is-4" v-if="itemSearch !== ''">No items found.</p>
                </div>
              </div>
            </div>
          </div>
          <br/>
        </div>
        <p class="subtitle is-4">
          <b>Total items</b>: {{ obtainedCount }}/{{ totalItems }}
          <br/>
          <b>Total price</b>: ${{ currentPrice }}/${{ totalPrice }} ({{ Math.round(currentPrice / totalPrice * 100 * 100) / 100 || 0 }}%)
        </p>
        <b-button type="is-success" size="is-medium" rounded @click="PatchShoppingListCompleted(id, !completed)">Mark as {{ completed === false ? 'completed' : 'uncompleted' }}</b-button>
        <b-button type="is-danger" size="is-medium" rounded @click="DeleteShoppingList(id)">Delete list</b-button>
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
      notesFromEmpty: false,
      itemSearch: '',
      authorNames: '',
      authorLastNames: '',
      totalItems: 0,
      loopCreated: new Date(),
      sortBy: shoppinglistCommon.GetShoppingListSortBy(),
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
    itemCard: () => import('@/frontend/components/authenticated/shopping-list-item-card-view.vue')
  },
  computed: {
    listItems () {
      var vm = this
      return this.RestructureShoppingListToTags(this.list.filter((item) => {
        return item.name.toLowerCase().indexOf(vm.itemSearch.toLowerCase()) !== -1
      }))
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
          window.clearInterval(this.intervalLoop)
          shoppinglist.DeleteShoppingList(id).then(resp => {
            common.DisplaySuccessToast('Deleted the shopping list')
            shoppinglistCommon.DeleteShoppingListFromCache(id)
            setTimeout(() => {
              this.$router.push({ name: 'Shopping list' })
            }, 1 * 1000)
          }).catch(err => {
            common.DisplayFailureToast('Failed to delete the shopping list' + '<br/>' + err.response.data.metadata.response)
          })
        }
      })
    },
    GetShoppingListItems () {
      var sortBy
      if (this.sortBy === 'price') {
        sortBy = this.sortBy
      }
      shoppinglist.GetShoppingListItems(this.id, sortBy).then(resp => {
        var responseList = resp.data.list
        this.totalItems = responseList === null ? 0 : responseList.length
        if (this.list === null) {
          this.list = []
        }

        if (responseList !== this.list) {
          this.list = responseList || []
          shoppinglistCommon.WriteShoppingListToCache(this.id, this.list)
        }
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
      this.intervalLoop = window.setInterval(() => {
        if (shoppinglistCommon.GetShoppingListAutoRefresh() !== true) {
          this.LoopStop()
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
    }
  },
  watch: {
    sortBy () {
      shoppinglistCommon.WriteShoppingListSortBy(this.sortBy)
    }
  },
  async beforeMount () {
    this.list = shoppinglistCommon.GetShoppingListFromCache(this.id) || []
    this.GetShoppingList()
    this.GetShoppingListItems()
  },
  async created () {
    this.LoopStart()
    window.addEventListener('focus', () => {
      this.loopCreated = new Date()
    })
  },
  beforeDestroy () {
    this.LoopStop()
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
</style>
