<template>
  <div>
    <section>
      <div class="card pointer-cursor-on-hover" @click="goToRef('/apps/shopping-list/list/' + list.id)">
        <div class="card-content">
          <div class="media">
            <div class="media-left">
              <b-icon
                icon="cart-outline"
                :type="list.completed === true ? 'is-success' : ''"
                size="is-medium">
              </b-icon>
            </div>
            <div class="media-content">
              <div class="display-items-on-the-same-line">
                <p class="title is-4">{{ list.name }}</p>
              </div>
              <p class="subtitle is-6">
                <span v-if="list.creationTimestamp == list.modificationTimestamp">
                  Created {{ TimestampToCalendar(list.creationTimestamp) }}, by {{ authorNames }}
                </span>
                <span v-else>
                  Updated {{ TimestampToCalendar(list.modificationTimestamp) }}, by {{ authorLastNames }}
                </span>
              </p>
            </div>
            <div class="media-right">
              <b-icon icon="chevron-right" size="is-medium" type="is-midgray"></b-icon>
            </div>
          </div>
          <div class="content">
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
import common from '@/frontend/common/common'
import flatmates from '@/frontend/requests/authenticated/flatmates'

export default {
  name: 'shopping-list-card-view',
  props: {
    list: Object
  },
  data () {
    return {
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
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    }
  },
  async beforeMount () {
    var list = this.list
    var userId = list.author
    if (list.creationTimestamp !== list.modificationTimestamp) {
      userId = list.authorLast
    }
    flatmates.GetFlatmate(userId).then(resp => {
      this.authorNames = resp.data.spec.names
      if (list.creationTimestamp !== list.modificationTimestamp) {
        this.authorLastNames = resp.data.spec.names
      }
    }).catch(err => {
      common.DisplayFailureToast('Unable to find author of list' + '<br/>' + err.response.data.metadata.response)
    })
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
