<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><router-link :to="'/apps/shopping-list/list/' + shoppingListId">{{ shoppingListName }}</router-link></li>
              <li class="is-active"><router-link :to="'/apps/shopping-list/list/' + shoppingListId + '/item/' + id">{{ name || 'Unnamed item' }}</router-link></li>
            </ul>
        </nav>
        <div>
          <h1 class="title is-1">{{ name || 'Unnamed item' }}</h1>
          <p class="subtitle is-3">View or edit this item</p>
          <p>
            Added {{ TimestampToCalendar(creationTimestamp) }}, by <router-link tag="a" :to="'/apps/flatmates?id=' + author"> {{ authorNames }} </router-link>
          </p>
          <p v-if="creationTimestamp !== modificationTimestamp">
            Last updated {{ TimestampToCalendar(modificationTimestamp) }}, by <router-link tag="a" :to="'/apps/flatmates?id=' + authorLast"> {{ authorLastNames }} </router-link>
          </p>
          <br/>
          <b-field label="Name">
            <b-input
              type="text"
              v-model="name"
              size="is-medium"
              maxlength="30"
              icon="text"
              required>
            </b-input>
          </b-field>
          <b-field label="Notes">
            <b-input
              type="text"
              v-model="notes"
              size="is-medium"
              maxlength="40"
              icon="text">
            </b-input>
          </b-field>
          <b-field label="Price">
            <b-input
              type="number"
              step="0.01"
              placeholder="0.00"
              v-model="price"
              icon="currency-usd"
              size="is-medium">
            </b-input>
          </b-field>
          <b-field label="Quantity">
            <b-numberinput
              v-model="quantity"
              size="is-medium"
              icon="numeric">
            </b-numberinput>
          </b-field>
          <div>
            <label class="label">Tag</label>
            <b-field>
              <p class="control">
                <b-input
                  type="text"
                  v-model="tag"
                  icon="tag"
                  size="is-medium">
                </b-input>
              </p>
              <p class="control" v-if="tags.length > 0">
                <b-dropdown>
                  <button class="button is-primary" slot="trigger">
                    <b-icon icon="menu-down"></b-icon>
                  </button>

                  <b-dropdown-item v-for="existingTag in tags" v-bind:key="existingTag" :value="existingTag" @click="tag = existingTag">{{ existingTag }}</b-dropdown-item>
                </b-dropdown>
              </p>
            </b-field>
          </div>
          <br/>
          <br/>
          <b-button type="is-success" size="is-medium" rounded native-type="submit" @click="UpdateShoppingListItem(shoppingListId, id, name, notes, price, quantity, tag)">Update item</b-button>
          <b-button type="is-danger" size="is-medium" rounded native-type="submit" @click="DeleteShoppingListItem(shoppingListId, id)">Delete item</b-button>
        </div>
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
  name: 'shopping-item-view',
  data () {
    return {
      shoppingListId: this.$route.params.listId,
      shoppingListName: '',
      authorNames: '',
      authorLastNames: '',
      tags: [],
      id: this.$route.params.itemId,
      name: '',
      notes: '',
      price: 0,
      quantity: 1,
      tag: undefined,
      author: '',
      authorLast: '',
      creationTimestamp: 0,
      modificationTimestamp: 0
    }
  },
  methods: {
    UpdateShoppingListItem (listId, itemId, name, notes, price, quantity, tag, obtained) {
      if (notes === '') {
        notes = undefined
      }
      if (price === 0) {
        price = undefined
      }

      price = Number(price)
      shoppinglist.UpdateShoppingListItem(listId, itemId, name, notes, price, quantity, tag, obtained).then(resp => {
        var item = resp.data.spec
        if (item.id !== '' && typeof item.id !== 'undefined') {
          common.DisplaySuccessToast('Updated item successfully')
          this.$router.push({ path: '/apps/shopping-list/list/' + this.shoppingListId })
        } else {
          common.DisplayFailureToast('Unable to find created shopping item')
        }
      }).catch(err => {
        common.DisplayFailureToast('Failed to add shopping list item' + ' - ' + err.response.data.metadata.response)
      })
    },
    DeleteShoppingListItem (listId, itemId) {
      Dialog.confirm({
        title: 'Delete item',
        message: 'Are you sure that you wish to delete this shopping list item?' + '<br/>' + 'This action cannot be undone.',
        confirmText: 'Delete item',
        type: 'is-danger',
        hasIcon: true,
        onConfirm: () => {
          shoppinglist.DeleteShoppingListItem(listId, itemId).then(resp => {
            common.DisplaySuccessToast(resp.data.metadata.response)
            setTimeout(() => {
              this.$router.push({ path: '/apps/shopping-list/list/' + this.shoppingListId })
            }, 1 * 1000)
          }).catch(err => {
            common.DisplayFailureToast('Failed to delete shopping list item' + ' - ' + err.response.data.metadata.response)
          })
        }
      })
    },
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    }
  },
  async beforeMount () {
    shoppinglist.GetShoppingList(this.shoppingListId).then(resp => {
      var list = resp.data.spec
      this.shoppingListName = list.name
      return shoppinglist.GetShoppingListItem(this.shoppingListId, this.id)
    }).then(resp => {
      var item = resp.data.spec
      this.name = item.name
      this.notes = item.notes
      this.price = item.price
      this.quantity = item.quantity
      this.tag = item.tag
      this.obtained = item.obtained
      this.author = item.author
      this.authorLast = item.authorLast
      this.creationTimestamp = item.creationTimestamp
      this.modificationTimestamp = item.modificationTimestamp
      return flatmates.GetFlatmate(item.author)
    }).then(resp => {
      this.authorNames = resp.data.spec.names
      return flatmates.GetFlatmate(this.authorLast)
    }).then(resp => {
      this.authorLastNames = resp.data.spec.names
      return shoppinglist.GetShoppingListItemTags()
    }).then(resp => {
      this.tags = resp.data.list || []
    })
  }
}
</script>
