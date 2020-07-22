<template>
  <div>
    <div class="card">
      <div class="card-content card-content-list">
        <div class="media">
          <div class="media-left" @click="PatchItemObtained(item.id, !item.obtained)">
            <b-checkbox size="is-medium" v-model="item.obtained"></b-checkbox>
          </div>
          <div class="media-content pointer-cursor-on-hover" @click="$router.push({ name: 'View shopping list item', params: { listId: shoppingListId, itemId: id } })">
            <div class="block">
              <p :class="item.obtained === true ? 'obtained' : ''" class="subtitle is-4 is-marginless">
                {{ item.name }}
                <span v-if="typeof item.price !== 'undefined' && item.price !== 0"> (${{ item.price.toFixed(2) }}) </span>
                <b v-if="item.quantity > 1">x{{ item.quantity }} </b>
                <b-icon
                  v-if="typeof item.price === 'undefined' || item.price === 0"
                  icon="currency-usd-off"
                  type="is-lightred"
                  size="is-small">
                </b-icon>
              </p>
              <span>
                <p class="subtitle is-6">
                  <span v-if="displayTag === true">
                    {{ item.tag }}
                  </span>
                  <span v-if="displayTag === true && typeof item.tag !== 'undefined' && typeof item.notes !== 'undefined' && item.notes !== ''">
                    -
                  </span>
                  <i>
                    {{ item.notes }}
                  </i>
                </p>
              </span>
            </div>
          </div>
          <div class="media-right">
            <b-tooltip label="Duplicate" class="is-paddingless" :delay="200">
              <b-button
                type="is-white"
                icon-right="content-duplicate"
                v-if="deviceIsMobile === false"
                @click="PostShoppingListItem(listId, item.name, item.notes, item.price, item.quantity, item.tag)" />
            </b-tooltip>

            <b-tooltip label="Delete" class="is-paddingless" :delay="200">
              <b-button
                type="is-danger"
                icon-right="delete"
                :loading="itemDeleting"
                v-if="deviceIsMobile === false"
                @click="DeleteShoppingListItem(item.id, index)" />
            </b-tooltip>
            <span class="pointer-cursor-on-hover" @click="$router.push({ name: 'View shopping list item', params: { listId: shoppingListId, itemId: id } })">
                <b-icon icon="chevron-right" size="is-medium" type="is-midgray"></b-icon>
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import common from '@/frontend/common/common'
import { DialogProgrammatic as Dialog } from 'buefy'
import shoppinglist from '@/frontend/requests/authenticated/shoppinglist'

export default {
  name: 'shopping list item card view',
  data () {
    return {
      itemDeleting: false
    }
  },
  props: {
    item: Object,
    listId: String,
    list: Object,
    index: Number,
    displayTag: Boolean,
    deviceIsMobile: Boolean,
    itemDisplayState: Number
  },
  methods: {
    PatchItemObtained (itemId, obtained) {
      shoppinglist.PatchShoppingListItemObtained(this.listId, itemId, obtained).then(() => {
        var displayAll = typeof this.itemDisplayState === 'number' && this.itemDisplayState === 0
        if (displayAll === true) {
          return
        }
        this.list.splice(this.index, 1)
      }).catch(err => {
        common.DisplayFailureToast('Failed to patch the obtained field of this item' + '<br/>' + err.response.data.metadata.response)
      })
    },
    DeleteShoppingListItem (itemId, index) {
      Dialog.confirm({
        title: 'Delete item',
        message: 'Are you sure that you wish to delete this shopping list item?' + '<br/>' + 'This action cannot be undone.',
        confirmText: 'Delete item',
        type: 'is-danger',
        hasIcon: true,
        onConfirm: () => {
          this.itemDeleting = true
          shoppinglist.DeleteShoppingListItem(this.listId, itemId).then(resp => {
            common.DisplaySuccessToast(resp.data.metadata.response)
            this.list.splice(this.index, 1)
          }).catch(err => {
            common.DisplayFailureToast('Failed to delete shopping list item' + ' - ' + err.response.data.metadata.response)
            this.itemDeleting = false
          })
        }
      })
    },
    PostShoppingListItem (listId, name, notes, price, quantity, tag) {
      Dialog.confirm({
        title: 'Duplicate item',
        message: 'Are you sure that you wish to duplicate this shopping list item?',
        confirmText: 'Duplicate item',
        type: 'is-warning',
        hasIcon: true,
        onConfirm: () => {
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
            if (item.id === '' || typeof item.id === 'undefined') {
              this.submitLoading = false
              common.DisplayFailureToast('Unable to find created shopping item')
            }
          }).catch(err => {
            this.submitLoading = false
            common.DisplayFailureToast(`Failed to add shopping list item - ${err.response.data.metadata.response}`)
          })
        }
      })
    }
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
