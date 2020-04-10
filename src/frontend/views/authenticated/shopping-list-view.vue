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
              <h1 class="title is-1 display-is-editable" @click="editing = !editing">{{ name || 'Unnamed list' }}</h1>
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
          {{ TimestampToCalendar(creationTimestamp) }}, by <router-link tag="a" :to="'/apps/flatmates?id=' + author"> {{ authorNames }} </router-link>
          </p>
          <p v-if="author !== authorLast && creationTimestamp !== modificationTimestamp">
            Last updated {{ TimestampToCalendar(creationTimestamp) }}, by <router-link tag="a" :to="'/apps/flatmates?id=' + author"> {{ authorLastNames }} </router-link>
          </p>
          <br/>
        </div>
        <div v-if="notes != '' || notesFromEmpty || editing">
          <div v-if="editing">
            <b-field>
              <b-input size="is-medium" maxlength="100" type="textarea" v-model="notes" class="subtitle is-3"></b-input>
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
        <b-button type="is-text" @click="() => { notesFromEmpty = true; editing = true }" v-if="!editing && notes.length == 0">Add notes</b-button>
        <b-button type="is-info" @click="() => { notesFromEmpty = false; editing = false; PatchShoppingList(name, notes) }" v-if="editing">Done</b-button>
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
        <section v-for="itemTag in list" v-bind:key="itemTag">
          <p class="title is-5">{{ itemTag.tag }}</p>
          <div v-for="item in itemTag.items" v-bind:key="item">
            <div class="card">
              <div class="card-content card-content-list">
                <div class="media">
                  <div class="media-left" @click="PatchItemObtained(item.id, !item.obtained)">
                    <b-checkbox size="is-medium" v-model="item.obtained"></b-checkbox>
                  </div>
                  <div class="media-content pointer-cursor-on-hover" @click="goToRef('/apps/shopping-list/list/' + id + '/item/' + item.id)">
                    <div class="block">
                      <p :class="item.obtained === true ? 'obtained' : ''" class="subtitle is-4">
                        {{ item.name }}
                        <b v-if="item.quantity > 1">x{{ item.quantity }}</b>
                        <span v-if="typeof item.price !== 'undefined' || item.price !== 0"> (${{ item.price }}) </span>
                        <b-icon
                          v-if="item.notes.length > 0"
                          icon="note-text-outline"
                          size="is-small">
                        </b-icon>
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <br/>
        </section>
        <br/>
        <p>
          <b>Total items</b>: {{ obtainedCount }}/{{ totalItems }}
          <br/>
          <b>Total price</b>: ${{ totalPrice }}
        </p>
        <br/>
        <br/>
        <b-button type="is-info" @click="MarkListAsCompleted(id)">Mark as completed</b-button>
        <b-button type="is-danger" @click="DeleteShoppingList(id)">Delete list</b-button>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/frontend/common/common'
import shoppinglist from '@/frontend/requests/authenticated/shoppinglist'
import flatmates from '@/frontend/requests/authenticated/flatmates'
import { DialogProgrammatic as Dialog } from 'buefy'

export default {
  name: 'Shopping List',
  data () {
    return {
      editing: false,
      notesFromEmpty: false,
      authorNames: '',
      authorLastNames: '',
      totalItems: 0,
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
  computed: {
    obtainedCount () {
      var obtained = 0
      var list = this.list
      for (var itemTag in list) {
        for (var item in list[itemTag].items) {
          if (list[itemTag].items[item].obtained === true) {
            obtained += 1
          }
        }
      }
      return obtained
    },
    totalPrice () {
      var totalPrice = 0
      var list = this.list
      for (var itemTag in list) {
        for (var item in list[itemTag].items) {
          console.log(list[itemTag].items[item].price)
          totalPrice += list[itemTag].items[item].price * list[itemTag].items[item].quantity
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
    GetShoppingList () {
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
    PatchShoppingList (name, notes) {
      var id = this.id
      shoppinglist.PatchShoppingList(id, name, notes).catch(err => {
        common.DisplayFailureToast('Failed to update shopping list' + '<br/>' + err.response.data.metadata.response)
      })
    },
    DeleteShoppingList (id) {
      Dialog.confirm({
        title: 'Delete shopping list',
        message: 'Are you sure that you wish to delete this shopping list?',
        confirmText: 'Delete shopping list',
        type: 'is-danger',
        hasIcon: true,
        onConfirm: () => {
          shoppinglist.DeleteShoppingList(id).then(resp => {
            common.DisplaySuccessToast('Deleted the shopping list')
            setTimeout(() => {
              this.$router.push({ name: 'Shopping list' })
            }, 1 * 1000)
          }).catch(err => {
            common.DisplayFailureToast('Failed to delete the shopping list' + '<br/>' + err.response.data.metadata.response)
          })
        }
      })
    },
    PatchItemObtained (itemId, obtained) {
      shoppinglist.PatchShoppingListItemObtained(this.id, itemId, obtained).catch(err => {
        common.DisplayFailureToast('Failed to patch the obtained field of this item' + '<br/>' + err.response.data.metadata.response)
      })
    },
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    }
  },
  async beforeMount () {
    this.GetShoppingList()
    shoppinglist.GetShoppingListItems(this.id).then(resp => {
      var responseList = resp.data.list
      this.totalItems = responseList === null ? 0 : responseList.length
      if (this.list === null) {
        this.list = []
      }

      // restructure be be organised by tag
      var currentTag = ''
      for (var item in responseList) {
        if (currentTag !== responseList[item].tag) {
          currentTag = responseList[item].tag
          var newItem = {
            tag: currentTag || 'Untagged',
            items: [responseList[item]]
          }

          this.list = [...this.list, newItem]
        } else {
          var currentListPosition = this.list.length - 1
          var currentSubListItems = this.list[currentListPosition].items

          this.list[currentListPosition].items = [...currentSubListItems, responseList[item]]
        }
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
