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
  <div>
    <section>
      <div class="card pointer-cursor-on-hover">
        <div class="card-content">
          <div class="media">
            <div class="media-left" @click="goToRef('/apps/shopping-list/list/' + list.id)">
              <b-icon
                icon="cart-outline"
                :type="list.completed === true ? 'is-success' : ''"
                size="is-medium">
              </b-icon>
            </div>
            <div class="media-content" @click="goToRef('/apps/shopping-list/list/' + list.id)">
              <div class="display-items-on-the-same-line">
                <p class="title is-4">{{ list.name }}</p>
              </div>
              <p class="subtitle is-6">
                <span v-if="list.creationTimestamp == list.modificationTimestamp">
                  Created {{ TimestampToCalendar(list.creationTimestamp) }}
                </span>
                <span v-else>
                  Updated {{ TimestampToCalendar(list.modificationTimestamp) }}
                </span>
              </p>
            </div>
            <div class="media-right">
              <b-tooltip label="Delete" class="is-paddingless" :delay="200">
                <b-button
                  type="is-danger"
                  icon-right="delete"
                  :loading="itemDeleting"
                  v-if="deviceIsMobile === false"
                  @click="DeleteShoppingList(list.id)" />
              </b-tooltip>
              <b-icon icon="chevron-right" size="is-medium" type="is-midgray"></b-icon>
            </div>
          </div>
          <div class="content" @click="goToRef('/apps/shopping-list/list/' + list.id)">
            <div>
              <b-tag type="is-info" v-if="list.completed">Completed</b-tag>
              <b-tag type="is-warning" v-if="!list.completed">Uncompleted</b-tag>
            </div>
            <br/>
            <span v-if="list.notes !== '' && typeof list.notes !== 'undefined'">
              <i>
                {{ PreviewNotes(list.notes) }}
              </i>
              <br/>
              <br/>
            </span>
            <div v-if="typeof list.count !== 'undefined' && list.count > 0">
              {{ list.count }} item(s)
            </div>
            <div v-else>
              0 items
            </div>
          </div>
        </div>
      </div>
      <br/>
    </section>
  </div>
</template>

<script>
import { DialogProgrammatic as Dialog } from 'buefy'
import common from '@/common/common'
import shoppinglist from '@/requests/authenticated/shoppinglist'
import shoppinglistCommon from '@/common/shoppinglist'

export default {
  name: 'shopping-list-card-view',
  props: {
    deviceIsMobile: Boolean,
    index: Number,
    lists: Object,
    list: Object
  },
  data () {
    return {
      deleteLoading: false,
      authorNames: '',
      authorLastNames: ''
    }
  },
  methods: {
    goToRef (ref) {
      this.$router.push({ path: ref })
    },
    PreviewNotes (notes) {
      if (notes.length <= 35) {
        return notes
      }
      var notesBytes = notes.split('')
      var notesBytesValid = notesBytes.filter((value, index) => {
        if (index <= 35) {
          return value
        }
      })
      return notesBytesValid.join('') + '...'
    },
    DeleteShoppingList (id) {
      Dialog.confirm({
        title: 'Delete shopping list',
        message: 'Are you sure that you wish to delete this shopping list?' + '<br/>' + 'This action cannot be undone.',
        confirmText: 'Delete shopping list',
        type: 'is-danger',
        hasIcon: true,
        onConfirm: () => {
          this.deleteLoading = true
          window.clearInterval(this.intervalLoop)
          shoppinglist.DeleteShoppingList(id).then(resp => {
            this.lists.splice(this.index, 1)
            common.DisplaySuccessToast('Deleted the shopping list')
            shoppinglistCommon.DeleteShoppingListFromCache(id)
          }).catch(err => {
            this.deleteLoading = false
            common.DisplayFailureToast('Failed to delete the shopping list' + '<br/>' + err.response.data.metadata.response)
          })
        }
      })
    },
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    }
  },
  async beforeMount () {
  }
}
</script>

<style>
.display-items-on-the-same-line {
    display: flex;

}

.display-items-on-the-same-line div {
    margin-left: 10px;
}
</style>
