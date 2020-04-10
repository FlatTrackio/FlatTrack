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
          <b-field label="Name">
            <b-input type="text"
                     v-model="name"
                     maxlength="30"
                     required>
            </b-input>
          </b-field>
          <b-field label="Notes">
            <b-input type="textarea"
                     v-model="notes"
                     maxlength="100"
                     >
            </b-input>
          </b-field>
          <b-field label="Price">
            <b-numberinput v-model="price" size="is-medium">
            </b-numberinput>
          </b-field>
          <b-field label="Quantity">
            <b-numberinput v-model="quantity" size="is-medium">
            </b-numberinput>
          </b-field>
          <b-button type="is-success" size="is-medium" rounded native-type="submit" @click="UpdateShoppingListItem(shoppingListId, id, name, notes, price, regular)" disabled>Update</b-button>
          <b-button type="is-danger" size="is-medium" rounded native-type="submit" @click="DeleteShoppingListItem(shoppingListId, id)">Delete</b-button>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/frontend/common/common'
import shoppinglist from '@/frontend/requests/authenticated/shoppinglist'
import { DialogProgrammatic as Dialog } from 'buefy'

export default {
  name: 'shopping-item-view',
  data () {
    return {
      shoppingListId: this.$route.params.listId,
      shoppingListName: '',
      id: this.$route.params.itemId,
      name: '',
      notes: '',
      price: 0,
      quantity: 1
    }
  },
  methods: {
    UpdateShoppingListItem (listId, name, notes, price, regular) {
      if (notes === '') {
        notes = undefined
      }
      if (price === 0) {
        price = undefined
      }

      shoppinglist.UpdateShoppingListItem(listId, name, notes, price, regular).then(resp => {
        var item = resp.data.spec
        if (item.id !== '' || typeof item.id === 'undefined') {
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
        message: 'Are you sure that you wish to delete this shopping list item?',
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
    }
  },
  async beforeMount () {
    shoppinglist.GetShoppingList(this.shoppingListId).then(resp => {
      var list = resp.data.spec
      this.shoppingListName = list.name
    })
    shoppinglist.GetShoppingListItem(this.shoppingListId, this.id).then(resp => {
      var item = resp.data.spec
      console.log({ item })
      this.name = item.name
      this.notes = item.notes
      this.price = item.price
      this.quantity = item.quantity
    })
  }
}
</script>
