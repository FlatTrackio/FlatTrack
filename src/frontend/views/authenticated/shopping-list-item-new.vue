<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><router-link :to="'/apps/shopping-list/list/' + id">{{ shoppingListName }}</router-link></li>
              <li class="is-active"><router-link :to="'/apps/shopping-list/list/' + id + '/new'">New shopping item</router-link></li>
            </ul>
        </nav>
        <div>
          <h1 class="title is-1">New shopping item</h1>
          <p class="subtitle is-3">Add an item to the list</p>
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
          <div class="field">
            <b-checkbox v-model="regular" size="is-medium">Regular item</b-checkbox>
          </div>
          <b-button type="is-success" size="is-medium" rounded native-type="submit" @click="PostShoppingListItem(shoppingListId, name, notes, price, regular)">Add</b-button>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/frontend/common/common'
import shoppinglist from '@/frontend/requests/authenticated/shoppinglist'

export default {
  name: 'shopping-item-new',
  data () {
    return {
      shoppingListId: this.$route.params.id,
      shoppingListName: '',
      name: '',
      notes: '',
      price: 0,
      regular: false,
      quantity: 1
    }
  },
  methods: {
    PostShoppingListItem (listId, name, notes, price, regular) {
      if (notes === '') {
        notes = undefined
      }
      if (price === 0) {
        price = undefined
      }

      shoppinglist.PostShoppingListItem(listId, name, notes, price, regular).then(resp => {
        var item = resp.data.spec
        if (item.id !== '' || typeof item.id === 'undefined') {
          this.$router.push({ path: '/apps/shopping-list/list/' + this.shoppingListId })
        } else {
          common.DisplayFailureToast('Unable to find created shopping item')
        }
      }).catch(err => {
        common.DisplayFailureToast(`Failed to add shopping list item - ${err.response.data.metadata.response}`)
      })
    }
  },
  async beforeMount () {
    shoppinglist.GetShoppingList(this.shoppingListId).then(resp => {
      var list = resp.data.spec
      this.shoppingListName = list.name
    })
  }
}
</script>
