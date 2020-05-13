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
          <b-field label="Name" class="is-marginless">
            <b-input
              type="text"
              v-model="name"
              size="is-medium"
              maxlength="30"
              icon="text"
              placeholder="Enter a name for this item"
              autofocus
              required>
            </b-input>
          </b-field>
          <b-field label="Notes (optional)" class="is-marginless">
            <b-input
              type="text"
              v-model="notes"
              size="is-medium"
              icon="text"
              placeholder="Enter extra information as notes to this item"
              maxlength="40">
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
          <b-field label="Quantity (optional)">
            <b-numberinput
              v-model="quantity"
              size="is-medium"
              placeholder="Enter how many of this item should be obtained"
              expanded
              min="1"
              controls-position="compact"
              icon="numeric">
            </b-numberinput>
          </b-field>
          <div>
            <label class="label">Tag (optional)</label>
            <b-field class="is-marginless">
              <p class="control" v-if="tags.length > 0">
                <b-dropdown>
                  <b-button
                    icon-left="menu-down"
                    type="is-primary"
                    slot="trigger"
                    size="is-medium">
                  </b-button>

                  <div v-for="existingTag in tags" v-bind:key="existingTag">
                    <b-dropdown-item v-if="existingTag !== '' && existingTag.length > 0 && typeof existingTag !== 'undefined'" :value="existingTag" @click="tag = existingTag">{{ existingTag }}</b-dropdown-item>
                  </div>
                </b-dropdown>
              </p>
              <b-input
                expanded
                type="text"
                v-model="tag"
                icon="tag"
                placeholder="Enter a tag to group the item"
                size="is-medium">
              </b-input>
            </b-field>
          </div>
          <br/>
          <b-button
            type="is-success"
            size="is-medium"
            icon-left="plus"
            native-type="submit"
            expanded
            :loading="submitLoading"
            @click="PostShoppingListItem(shoppingListId, name, notes, price, quantity, tag)">
            Add
          </b-button>
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
      submitLoading: false,
      name: '',
      notes: '',
      price: 0,
      quantity: 1,
      tag: undefined
    }
  },
  methods: {
    PostShoppingListItem (listId, name, notes, price, quantity, tag) {
      this.submitLoading = true
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
          this.submitLoading = false
          common.DisplayFailureToast('Unable to find created shopping item')
        }
      }).catch(err => {
        this.submitLoading = false
        common.DisplayFailureToast(`Failed to add shopping list item - ${err.response.data.metadata.response}`)
      })
    }
  },
  async beforeMount () {
    shoppinglist.GetShoppingList(this.shoppingListId).then(resp => {
      var list = resp.data.spec
      this.shoppingListName = list.name
      return shoppinglist.GetAllShoppingListItemTags()
    }).then(resp => {
      this.tags = resp.data.list || []
    })
  }
}
</script>
