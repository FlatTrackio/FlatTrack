<!--
     This program is free software: you can redistribute it and/or modify
     it under the terms of the Affero GNU General Public License as published by
     the Free Software Foundation, either version 3 of the License, or
     (at your option) any later version.

     This program is distributed in the hope that it will be useful,
     but WITHOUT ANY WARRANTY; without even the implied warranty of
     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
     GNU General Public License for more details.

     You should have received a copy of the Affero GNU General Public License
     along with this program.  If not, see <https://www.gnu.org/licenses/>.
-->

<template>
  <div class="item-page">
    <div class="modal-card" style="width: auto">
      <header class="modal-card-head">
        <p class="modal-card-title">New shopping item</p>
        <p class="modal-card-subtitle">Add an item to the list</p>
      </header>
      <section class="modal-card-body">
        <div>
          <b-field label="Name" class="is-marginless">
            <b-input
              type="text"
              v-model="name"
              size="is-medium"
              maxlength="30"
              icon="text"
              placeholder="Enter a name for this item"
              icon-right="close-circle"
              icon-right-clickable
              @icon-right-click="name = ''"
              @keyup.enter.native="PostShoppingListItem"
              autofocus
              required
            >
            </b-input>
          </b-field>
          <b-field label="Notes (optional)" class="is-marginless">
            <b-input
              type="text"
              v-model="notes"
              size="is-medium"
              icon="text"
              placeholder="Enter information extra"
              @keyup.enter.native="PostShoppingListItem"
              icon-right="close-circle"
              icon-right-clickable
              @icon-right-click="notes = ''"
              maxlength="40"
            >
            </b-input>
          </b-field>
          <b-field label="Price (optional)">
            <b-input
              type="number"
              step="0.01"
              placeholder="0.00"
              v-model="price"
              icon="currency-usd"
              icon-right="close-circle"
              icon-right-clickable
              @icon-right-click="price = ''"
              @keyup.enter.native="PostShoppingListItem"
              size="is-medium"
            >
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
              icon="numeric"
            >
            </b-numberinput>
          </b-field>
          <div>
            <div class="field has-addons">
              <label class="label">Tag (optional)</label>
              <p class="control">
                <infotooltip
                  message="To manage tags, navigate to the Apps -> Shopping List -> Manage tags page"
                />
              </p>
            </div>
            <b-field class="is-marginless">
              <b-dropdown>
                <b-button
                  icon-left="menu-down"
                  type="is-primary"
                  slot="trigger"
                  size="is-medium"
                  v-if="tagsList.length > 0 || tags.length > 0"
                >
                </b-button>

                <b-dropdown-item disabled v-if="tagsList.length > 0"
                  >Tags in this list</b-dropdown-item
                >
                <div
                  v-for="existingListTag in tagsList"
                  v-bind:key="existingListTag"
                >
                  <b-dropdown-item
                    v-if="
                      existingListTag !== '' &&
                      existingListTag.length > 0 &&
                      typeof existingListTag !== 'undefined'
                    "
                    :value="existingListTag"
                    @click="tag = existingListTag"
                    >{{ existingListTag }}</b-dropdown-item
                  >
                </div>
                <b-dropdown-item disabled v-if="tags.length > 0"
                  >Tags in all lists</b-dropdown-item
                >
                <div v-for="existingTag in tags" v-bind:key="existingTag">
                  <b-dropdown-item
                    v-if="
                      existingTag.name !== '' &&
                      existingTag.name.length > 0 &&
                      typeof existingTag.name !== 'undefined'
                    "
                    :value="existingTag.name"
                    @click="tag = existingTag.name"
                    >{{ existingTag.name }}</b-dropdown-item
                  >
                </div>
              </b-dropdown>
              <b-input
                expanded
                type="text"
                v-model="tag"
                icon="tag"
                maxlength="30"
                placeholder="Enter a tag to group the item"
                @keyup.enter.native="PostShoppingListItem"
                icon-right="close-circle"
                icon-right-clickable
                @icon-right-click="tag = ''"
                size="is-medium"
              >
              </b-input>
            </b-field>
          </div>
          <b-field label="Obtained">
            <b-checkbox size="is-medium" v-model="obtained">
              Obtained
            </b-checkbox>
          </b-field>
          <p
            v-if="typeof price !== 'undefined' && price !== 0 && quantity > 1"
            class="m-1"
          >
            Total price with quantity: ${{ itemCurrentPrice.toFixed(2) }}
          </p>
          <div class="level">
            <b-button
              type="is-warning"
              size="is-medium"
              icon-left="arrow-left"
              native-type="submit"
              @click="$parent.close()"
            >
              Back
            </b-button>
            <b-button
              type="is-success"
              size="is-medium"
              icon-left="plus"
              native-type="submit"
              expanded
              :loading="submitLoading"
              :disabled="submitLoading"
              @click="PostShoppingListItem"
            >
              Add item
            </b-button>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/common/common'
import shoppinglist from '@/requests/authenticated/shoppinglist'

export default {
  name: 'shopping-item-new',
  components: {
    infotooltip: () => import('@/components/common/info-tooltip.vue')
  },
  data () {
    return {
      shoppingListId: this.$route.params.id,
      shoppingListName: '',
      tags: [],
      tagsList: [],
      submitLoading: false,
      name: '',
      notes: '',
      price: 0,
      quantity: 1,
      tag: '',
      obtained: false
    }
  },
  props: {
    withTag: String,
    withName: String
  },
  methods: {
    CopyHrefToClipboard () {
      common.CopyHrefToClipboard()
    },
    PostShoppingListItem () {
      this.submitLoading = true
      if (this.notes === '') {
        this.notes = undefined
      }
      if (this.price === 0) {
        this.price = undefined
      } else {
        this.price = parseFloat(this.price)
      }
      if (this.tag === '') {
        this.tag = 'Untagged'
      }

      shoppinglist
        .PostShoppingListItem(
          this.shoppingListId,
          this.name,
          this.notes,
          this.price,
          this.quantity,
          this.tag,
          this.obtained
        )
        .then((resp) => {
          var item = resp.data.spec
          if (item.id !== '' || typeof item.id === 'undefined') {
            this.$parent.close()
          } else {
            this.submitLoading = false
            common.DisplayFailureToast('Unable to find created shopping item')
          }
        })
        .catch((err) => {
          this.submitLoading = false
          common.DisplayFailureToast(
            `Failed to add shopping list item - ${err.response.data.metadata.response}`
          )
        })
    }
  },
  computed: {
    itemCurrentPrice () {
      return this.price * this.quantity
    }
  },
  async beforeMount () {
    if (this.withTag && this.tag === '') {
      this.tag = this.withTag
    }
    console.log(this.tag, this.tag === '', this.withTag)
    if (this.withName && this.name === '') {
      this.name = this.withName
    }
    console.log(this.tag, this.withTag)
    shoppinglist
      .GetShoppingList(this.shoppingListId)
      .then((resp) => {
        var list = resp.data.spec
        this.shoppingListName = list.name
        return shoppinglist.GetAllShoppingListItemTags()
      })
      .then((resp) => {
        this.tags = resp.data.list || []
        return shoppinglist.GetShoppingListItemTags(this.shoppingListId)
      })
      .then((resp) => {
        this.tagsList = resp.data.list || []
      })
  }
}
</script>
