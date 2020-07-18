<template>
  <div class="card pointer-cursor-on-hover">
    <div class="card-content card-content-list">
      <div class="media">
        <div class="media-left" @click="UpdateTag(tag)">
          <b-icon icon="tag" size="is-medium" type="is-midgray"></b-icon>
        </div>
        <div class="media-content" @click="UpdateTag(tag)">
          <p class="title is-4"> {{ tag.name }} </p>
          <p class="subtitle is-6">
            <span v-if="tag.creationTimestamp == tag.modificationTimestamp">
              Created {{ TimestampToCalendar(tag.creationTimestamp) }}, by {{ authorNames }}
            </span>
            <span v-else>
              Updated {{ TimestampToCalendar(tag.modificationTimestamp) }}, by {{ authorNamesLast }}
            </span>
          </p>
        </div>
        <div class="media-right">
          <!-- Delete button -->
          <b-tooltip label="Delete" class="is-paddingless" :delay="200">
            <b-button
              type="is-danger"
              icon-right="delete"
              @click="DeleteShoppingListTag(tag.id, index)" />
          </b-tooltip>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import common from '@/frontend/common/common'
import shoppinglist from '@/frontend/requests/authenticated/shoppinglist'
import flatmates from '@/frontend/requests/authenticated/flatmates'
import { DialogProgrammatic as Dialog } from 'buefy'

export default {
  name: 'Shopping Tag card',
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
          shoppinglist.UpdateShoppingTag(tag.id, value).then(resp => {
            this.pageLoading = true
            tag.name = resp.data.spec.name
            common.DisplaySuccessToast(resp.data.metadata.response)
            this.$emit('displayFloatingAddButton', true)
          }).catch(err => {
            common.DisplayFailureToast(`Failed to create tag; ${err.response.data.metadata.response}`)
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
        message: 'Are you sure that you wish to delete this shopping list tag?' + '<br/>' + 'This action cannot be undone.' + '<br/>' + '<br/>' + 'Please note: this will not alter existing items with this tag',
        confirmText: 'Delete tag',
        type: 'is-danger',
        hasIcon: true,
        onConfirm: () => {
          this.itemDeleting = true
          shoppinglist.DeleteShoppingTag(id).then(resp => {
            common.DisplaySuccessToast(resp.data.metadata.response)
            this.list.splice(index, 1)
            this.$emit('displayFloatingAddButton', true)
          }).catch(err => {
            common.DisplayFailureToast('Failed to delete shopping tag' + ' - ' + err.response.data.metadata.response)
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
    flatmates.GetFlatmate(this.tag.author).then(resp => {
      this.authorNames = resp.data.spec.names
      return flatmates.GetFlatmate(this.tag.authorLast)
    }).then(resp => {
      this.authorNamesLast = resp.data.spec.names
    }).catch(err => {
      common.DisplayFailureToast('Unable to find author of tag' + '<br/>' + err.response.data.metadata.response)
    })
  }
}
</script>

<style scoped>

</style>
