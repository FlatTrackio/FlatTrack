<template>
  <div>
    <div class="card">
      <div class="card-content card-content-list">
        <div class="media">
          <div class="media-left" @click="PatchItemObtained(item.id, !item.obtained)">
            <b-checkbox size="is-medium" v-model="item.obtained"></b-checkbox>
          </div>
          <div class="media-content pointer-cursor-on-hover" @click="goToRef('/apps/shopping-list/list/' + listId + '/item/' + item.id)">
            <div class="block">
              <p :class="item.obtained === true ? 'obtained' : ''" class="subtitle is-4 is-marginless">
                {{ item.name }}
                <span v-if="typeof item.price !== 'undefined' && item.price !== 0"> (${{ item.price }}) </span>
                <b v-if="item.quantity > 1">x{{ item.quantity }} </b>
                <b-icon
                  v-if="typeof item.price === 'undefined' || item.price === 0"
                  icon="currency-usd-off"
                  size="is-small">
                </b-icon>
                <b-icon
                  v-if="item.notes.length > 0"
                  icon="note-text-outline"
                  size="is-small">
                </b-icon>
              </p>
              <span>
                <p class="subtitle is-6">
                  <span v-if="displayTag === true">
                    {{ item.tag }}
                  </span>
                  <span v-if="displayTag === true && typeof item.notes !== 'undefined' && item.notes !== ''">
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
            <b-button type="is-danger" icon-right="delete" v-if="deviceIsMobile === false" @click="DeleteShoppingListItem(listId, item.id, index)" />
            <b-icon icon="chevron-right" size="is-medium" type="is-midgray"></b-icon>
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
      deviceIsMobile: false
    }
  },
  props: {
    item: Object,
    listId: String,
    list: Object,
    index: Number,
    displayTag: Boolean
  },
  methods: {
    goToRef (ref) {
      this.$router.push({ path: ref })
    },
    PatchItemObtained (itemId, obtained) {
      shoppinglist.PatchShoppingListItemObtained(this.listId, itemId, obtained).catch(err => {
        common.DisplayFailureToast('Failed to patch the obtained field of this item' + '<br/>' + err.response.data.metadata.response)
      })
    },
    AdjustForMobile () {
      this.deviceIsMobile = common.DeviceIsMobile()
    },
    DeleteShoppingListItem (itemId, index) {
      Dialog.confirm({
        title: 'Delete item',
        message: 'Are you sure that you wish to delete this shopping list item?' + '<br/>' + 'This action cannot be undone.',
        confirmText: 'Delete item',
        type: 'is-danger',
        hasIcon: true,
        onConfirm: () => {
          shoppinglist.DeleteShoppingListItem(this.listId, itemId).then(resp => {
            common.DisplaySuccessToast(resp.data.metadata.response)
            this.list.splice(index, 1)
          }).catch(err => {
            common.DisplayFailureToast('Failed to delete shopping list item' + ' - ' + err.response.data.metadata.response)
          })
        }
      })
    }
  },
  async created () {
    this.AdjustForMobile()
    window.addEventListener('resize', this.AdjustForMobile.bind(this))
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
