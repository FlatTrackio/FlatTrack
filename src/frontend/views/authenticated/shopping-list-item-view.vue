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
          <b-loading :is-full-page="false" :active.sync="itemIsLoading" :can-cancel="false"></b-loading>
          <h1 class="title is-1">{{ name || 'Unnamed item' }}</h1>
          <p class="subtitle is-4">View or edit this item</p>
          <b-field label="Name">
            <b-input
              type="text"
              v-model="name"
              size="is-medium"
              maxlength="30"
              icon="text"
              placeholder="Enter a name for this item"
              required>
            </b-input>
          </b-field>
          <b-field label="Notes (optional)">
            <b-input
              type="text"
              v-model="notes"
              size="is-medium"
              maxlength="40"
              placeholder="Enter extra information as notes to this item"
              icon="text">
            </b-input>
          </b-field>
          <b-field label="Price (optional)">
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
              placeholder="Enter how many of this item should be obtained"
              min="0"
              expanded
              required
              controls-position="compact"
              icon="numeric">
            </b-numberinput>
          </b-field>
          <div>
            <div class="field has-addons">
              <label class="label">Tag (optional)</label>
              <p class="control">
                <infotooltip message="To manage tags, navigate to the Apps -> Shopping List -> Manage tags page"/>
              </p>
            </div>
            <b-field>
              <p class="control" v-if="tags.length > 0">
                <b-dropdown>
                  <b-button
                    icon-left="menu-down"
                    type="is-primary"
                    slot="trigger"
                    size="is-medium">
                  </b-button>

                  <div v-for="existingTag in tags" v-bind:key="existingTag">
                    <b-dropdown-item v-if="existingTag.name !== '' && existingTag.name.length > 0 && typeof existingTag.name !== 'undefined'" :value="existingTag.name" @click="tag = existingTag.name">{{ existingTag.name }}</b-dropdown-item>
                  </div>
                </b-dropdown>
              </p>
              <b-input
                expanded
                type="text"
                v-model="tag"
                icon="tag"
                maxlength="30"
                placeholder="Enter a tag to group the item"
                size="is-medium">
              </b-input>
            </b-field>
          </div>
          <br/>
          <b-field
            label="Obtained">
            <b-checkbox
              size="is-medium"
              v-model="obtained">
              Obtained
            </b-checkbox>
          </b-field>
          <b-field>
            <b-button
              type="is-success"
              size="is-medium"
              icon-left="delta"
              native-type="submit"
              expanded
              :loading="submitLoading"
              @click="UpdateShoppingListItem(shoppingListId, id, name, notes, price, quantity, tag, obtained)">
              Update item
            </b-button>
            <p class="control">
              <b-button
                type="is-danger"
                size="is-medium"
                icon-left="delete"
                native-type="submit"
                :loading="deleteLoading"
                @click="DeleteShoppingListItem(shoppingListId, id)">
              </b-button>
            </p>
          </b-field>
          <p>
            Added {{ TimestampToCalendar(creationTimestamp) }}, by <router-link tag="a" :to="'/apps/flatmates?id=' + author"> {{ authorNames }} </router-link>
          </p>
          <p v-if="creationTimestamp !== modificationTimestamp">
            Last updated {{ TimestampToCalendar(modificationTimestamp) }}, by <router-link tag="a" :to="'/apps/flatmates?id=' + authorLast"> {{ authorLastNames }} </router-link>
          </p>
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
  components: {
    infotooltip: () => import('@/frontend/components/common/info-tooltip.vue')
  },
  data () {
    return {
      shoppingListId: this.$route.params.listId,
      shoppingListName: '',
      authorNames: '',
      authorLastNames: '',
      tags: [],
      itemIsLoading: true,
      submitLoading: false,
      deleteLoading: false,
      id: this.$route.params.itemId,
      name: '',
      notes: '',
      price: 0,
      quantity: 1,
      tag: undefined,
      obtained: false,
      author: '',
      authorLast: '',
      creationTimestamp: 0,
      modificationTimestamp: 0
    }
  },
  methods: {
    UpdateShoppingListItem (listId, itemId, name, notes, price, quantity, tag, obtained) {
      this.submitLoading = true
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
          this.submitLoading = false
          common.DisplayFailureToast('Unable to find created shopping item')
        }
      }).catch(err => {
        common.DisplayFailureToast('Failed to add shopping list item' + ' - ' + err.response.data.metadata.response)
        this.submitLoading = false
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
          this.deleteLoading = true
          shoppinglist.DeleteShoppingListItem(listId, itemId).then(resp => {
            common.DisplaySuccessToast(resp.data.metadata.response)
            setTimeout(() => {
              this.$router.push({ path: '/apps/shopping-list/list/' + this.shoppingListId })
            }, 1 * 1000)
          }).catch(err => {
            this.deleteLoading = false
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
      return shoppinglist.GetAllShoppingListItemTags()
    }).then(resp => {
      this.itemIsLoading = false
      this.tags = resp.data.list || []
    }).catch(err => {
      if (err.response.status === 404) {
        common.DisplayFailureToast('Error item not found' + '<br/>' + err.response.data.metadata.response)
        this.$router.push({ path: '/apps/shopping-list/list/' + this.shoppingListId })
        return
      }
      common.DisplayFailureToast('Error loading the shopping list item' + '<br/>' + err.response.data.metadata.response)
    })
  }
}
</script>
