<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><router-link :to="'/apps/shopping-list/list/' + shoppingListId">{{ shoppingListName }}</router-link></li>
              <li class="is-active"><router-link :to="'/apps/shopping-list/list/' + shoppingListId + '/new'">New shopping item</router-link></li>
            </ul>
        </nav>
        <div>
          <h1 class="title is-1">New shopping item</h1>
          <p class="subtitle is-3">Add an item to the list</p>
          <b-field label="Name">
            <b-input
              type="text"
              v-model="name"
              size="is-medium"
              maxlength="30"
              icon="email"
              required>
            </b-input>
          </b-field>
          <b-field label="Notes">
            <b-input
              type="text"
              v-model="notes"
              size="is-medium"
              icon="text"
              maxlength="40">
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
          <b-button type="is-success" size="is-medium" rounded native-type="submit" @click="PostShoppingListItem(shoppingListId, name, notes, price, quantity, tag)">Add</b-button>
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
      tags: [],
      name: '',
      notes: '',
      price: 0,
      quantity: 1,
      tag: undefined
    }
  },
  methods: {
    PostShoppingListItem (listId, name, notes, price, quantity, tag) {
      if (notes === '') {
        notes = undefined
      }
      if (price === 0) {
        price = undefined
      } else {
        price = parseFloat(price)
      }

      shoppinglist.PostShoppingListItem(listId, name, notes, price, quantity, tag).then(resp => {
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
      return shoppinglist.GetShoppingListItemTags()
    }).then(resp => {
      this.tags = resp.data.list || []
    })
  }
}
</script>
