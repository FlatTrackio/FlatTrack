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
        <p class="modal-card-title">{{ name || "Unnamed item" }}</p>
        <p class="modal-card-subtitle">View or edit this item</p>
      </header>
      <section class="modal-card-body">
        <div>
          <b-loading
            :is-full-page="false"
            :active.sync="itemIsLoading"
            :can-cancel="false"
          ></b-loading>
          <b-field label="Name">
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
              @keyup.enter.native="UpdateShoppingListItem"
              required
            >
            </b-input>
          </b-field>
          <b-field label="Notes (optional)">
            <b-input
              type="text"
              v-model="notes"
              size="is-medium"
              maxlength="40"
              placeholder="Enter extra information as notes to this item"
              icon-right="close-circle"
              icon-right-clickable
              @icon-right-click="notes = ''"
              @keyup.enter.native="UpdateShoppingListItem"
              icon="text"
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
              @keyup.enter.native="UpdateShoppingListItem"
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
              @keyup.enter.native="UpdateShoppingListItem"
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
            <b-field>
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
                icon-right="close-circle"
                icon-right-clickable
                @icon-right-click="tag = ''"
                @keyup.enter.native="UpdateShoppingListItem"
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
            class="pb-2"
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
              icon-left="delta"
              native-type="submit"
              expanded
              :loading="submitLoading"
              :disabled="submitLoading"
              @click="UpdateShoppingListItem"
            >
              Update item
            </b-button>
            <b-button
              type="is-danger"
              size="is-medium"
              icon-left="delete"
              native-type="submit"
              :loading="deleteLoading"
              @click="DeleteShoppingListItem(shoppingListId, id)"
            >
            </b-button>
          </div>
          <p>
            Added {{ TimestampToCalendar(creationTimestamp) }}, by
            <router-link
              tag="a"
              :to="{ name: 'My Flatmates', query: { id: author } }"
              >{{ authorNames }}</router-link
            >
            <span v-if="templateId">
              (templated from
              <router-link
                tag="a"
                :to="{ name: 'View shopping list', params: { id: templateId } }"
                >{{ templateListName }}</router-link
              >)
            </span>
          </p>
          <p v-if="creationTimestamp !== modificationTimestamp">
            Last updated {{ TimestampToCalendar(modificationTimestamp) }}, by
            <router-link
              tag="a"
              :to="{ name: 'My Flatmates', query: { id: authorLast } }"
              >{{ authorLastNames }}</router-link
            >
          </p>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/common/common'
import shoppinglist from '@/requests/authenticated/shoppinglist'
import flatmates from '@/requests/authenticated/flatmates'
import { DialogProgrammatic as Dialog } from 'buefy'

export default {
  name: 'shopping-item-view',
  components: {
    infotooltip: () => import('@/components/common/info-tooltip.vue')
  },
  data () {
    return {
      prevRoute: null,
      shoppingListName: '',
      authorNames: '',
      authorLastNames: '',
      tags: [],
      tagsList: [],
      itemIsLoading: true,
      submitLoading: false,
      deleteLoading: false,
      templateListName: '',
      name: '',
      notes: '',
      price: 0,
      quantity: 1,
      tag: undefined,
      obtained: false,
      author: '',
      authorLast: '',
      templateId: undefined,
      creationTimestamp: 0,
      modificationTimestamp: 0
    }
  },
  props: {
    id: String,
    shoppingListId: String
  },
  methods: {
    CopyHrefToClipboard () {
      common.CopyHrefToClipboard()
    },
    UpdateShoppingListItem () {
      this.submitLoading = true
      if (this.notes === '') {
        this.notes = undefined
      }
      if (this.price === 0) {
        this.price = undefined
      }
      if (this.tag === '') {
        this.tag = 'Untagged'
      }

      this.price = Number(this.price)
      shoppinglist
        .UpdateShoppingListItem(
          this.shoppingListId,
          this.id,
          this.name,
          this.notes,
          this.price,
          this.quantity,
          this.tag,
          this.obtained
        )
        .then((resp) => {
          var item = resp.data.spec
          if (item.id !== '' && typeof item.id !== 'undefined') {
            common.DisplaySuccessToast('Updated item successfully')
            this.$parent.close()
          } else {
            this.submitLoading = false
            common.DisplayFailureToast('Unable to find created shopping item')
          }
        })
        .catch((err) => {
          common.DisplayFailureToast(
            'Failed to add shopping list item' +
              ' - ' +
              err.response.data.metadata.response
          )
          this.submitLoading = false
        })
    },
    DeleteShoppingListItem (listId, itemId) {
      Dialog.confirm({
        title: 'Delete item',
        message:
          'Are you sure that you wish to delete this shopping list item?' +
          '<br/>' +
          'This action cannot be undone.',
        confirmText: 'Delete item',
        type: 'is-danger',
        hasIcon: true,
        onConfirm: () => {
          this.deleteLoading = true
          shoppinglist
            .DeleteShoppingListItem(listId, itemId)
            .then((resp) => {
              common.DisplaySuccessToast(resp.data.metadata.response)
              setTimeout(() => {
                this.$parent.close()
              }, 1 * 1000)
            })
            .catch((err) => {
              this.deleteLoading = false
              common.DisplayFailureToast(
                'Failed to delete shopping list item' +
                  ' - ' +
                  err.response.data.metadata.response
              )
            })
        }
      })
    },
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    }
  },
  computed: {
    itemCurrentPrice () {
      return this.price * this.quantity
    }
  },
  async beforeMount () {
    shoppinglist
      .GetShoppingList(this.shoppingListId)
      .then((resp) => {
        var list = resp.data.spec
        this.shoppingListName = list.name
        return shoppinglist.GetShoppingListItem(this.shoppingListId, this.id)
      })
      .then((resp) => {
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
        this.templateId = item.templateId
        return flatmates.GetFlatmate(item.author)
      })
      .then((resp) => {
        this.authorNames = resp.data.spec.names
        return flatmates.GetFlatmate(this.authorLast)
      })
      .then((resp) => {
        this.authorLastNames = resp.data.spec.names
        return shoppinglist.GetAllShoppingListItemTags()
      })
      .then((resp) => {
        this.itemIsLoading = false
        this.tags = resp.data.list || []
        return shoppinglist.GetShoppingListItemTags(this.shoppingListId)
      })
      .then((resp) => {
        this.tagsList = resp.data.list || []
        if (typeof this.templateId === 'undefined' || this.templateId === '') {
          return
        }
        return shoppinglist.GetShoppingList(this.templateId)
      })
      .then((resp) => {
        if (typeof this.templateId === 'undefined' || this.templateId === '') {
          return
        }
        this.templateListName = resp.data.spec.name
      })
      .catch((err) => {
        if (err.response.status === 404) {
          common.DisplayFailureToast(
            'Error item not found' +
              '<br/>' +
              err.response.data.metadata.response
          )
          this.$router.push({
            name: 'View shopping list',
            params: { id: this.shoppingListId }
          })
          return
        }
        common.DisplayFailureToast(
          'Error loading the shopping list item' +
            '<br/>' +
            err.response.data.metadata.response
        )
      })
  }
}
</script>
