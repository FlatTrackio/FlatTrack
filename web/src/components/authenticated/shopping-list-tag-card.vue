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
  <div class="card pointer-cursor-on-hover">
    <div class="card-content card-content-list">
      <div class="media">
        <div class="media-left" @click="UpdateTag(tag)">
          <b-icon icon="tag" size="is-medium" type="is-midgray"></b-icon>
        </div>
        <div class="media-content" @click="UpdateTag(tag)">
          <p class="title is-4">{{ tag.name }}</p>
          <p class="subtitle is-6">
            <span v-if="tag.creationTimestamp == tag.modificationTimestamp">
              Created {{ TimestampToCalendar(tag.creationTimestamp) }}, by
              {{ authorNames }}
            </span>
            <span v-else>
              Updated {{ TimestampToCalendar(tag.modificationTimestamp) }}, by
              {{ authorNamesLast }}
            </span>
          </p>
        </div>
        <div class="media-right">
          <!-- Delete button -->
          <b-tooltip label="Delete" class="is-paddingless" :delay="200">
            <b-button
              type="is-danger"
              icon-right="delete"
              @click="DeleteShoppingListTag(tag.id, index)"
            />
          </b-tooltip>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import common from '@/common/common'
import shoppinglist from '@/requests/authenticated/shoppinglist'
import flatmates from '@/requests/authenticated/flatmates'
import { DialogProgrammatic as Dialog } from 'buefy'

export default {
  name: 'shopping-tag-card',
  data () {
    return {
      authorNames: '...',
      authorNamesLast: '...'
    }
  },
  props: {
    tag: Object,
    index: Number,
    tags: Object,
    displayFloatingAddButton: Boolean
  },
  methods: {
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    },
    UpdateTag (tag) {
      this.$emit('displayFloatingAddButton', false)
      Dialog.prompt({
        title: 'Edit tag',
        message: `Enter the name of a tag to create.`,
        container: null,
        icon: 'tag',
        hasIcon: true,
        inputAttrs: {
          placeholder: 'e.g. Fruits and Veges',
          maxlength: 30,
          value: tag.name
        },
        trapFocus: true,
        onConfirm: (value) => {
          shoppinglist
            .UpdateShoppingTag(tag.id, value)
            .then((resp) => {
              this.pageLoading = true
              tag.name = resp.data.spec.name
              common.DisplaySuccessToast(resp.data.metadata.response)
              this.$emit('displayFloatingAddButton', true)
            })
            .catch((err) => {
              common.DisplayFailureToast(
                `Failed to create tag; ${err.response.data.metadata.response}`
              )
              this.$emit('displayFloatingAddButton', true)
            })
        },
        onCancel: () => {
          this.$emit('displayFloatingAddButton', true)
        }
      })
    },
    DeleteShoppingListTag (id, index) {
      this.$emit('displayFloatingAddButton', false)
      Dialog.confirm({
        title: 'Delete tag',
        message:
          'Are you sure that you wish to delete this shopping list tag?' +
          '<br/>' +
          'This action cannot be undone.' +
          '<br/>' +
          '<br/>' +
          'Please note: this will not alter existing items with this tag',
        confirmText: 'Delete tag',
        type: 'is-danger',
        hasIcon: true,
        onConfirm: () => {
          this.itemDeleting = true
          shoppinglist
            .DeleteShoppingTag(id)
            .then((resp) => {
              common.DisplaySuccessToast(resp.data.metadata.response)
              let removedFromTags = this.tags
              removedFromTags.splice(this.index, 1)
              this.$emit('tags', removedFromTags)
              this.$emit('displayFloatingAddButton', true)
            })
            .catch((err) => {
              common.DisplayFailureToast(
                'Failed to delete shopping tag' +
                  ' - ' +
                  err.response.data.metadata.response
              )
              this.itemDeleting = false
              this.$emit('displayFloatingAddButton', true)
            })
        },
        onCancel: () => {
          this.$emit('displayFloatingAddButton', true)
        }
      })
    }
  },
  async beforeMount () {
    flatmates
      .GetFlatmate(this.tag.author)
      .then((resp) => {
        this.authorNames = resp.data.spec.names
        return flatmates.GetFlatmate(this.tag.authorLast)
      })
      .then((resp) => {
        this.authorNamesLast = resp.data.spec.names
      })
      .catch((err) => {
        common.DisplayFailureToast(
          'Unable to find author of tag' +
            '<br/>' +
            err.response.data.metadata.response
        )
      })
  }
}
</script>

<style scoped></style>
