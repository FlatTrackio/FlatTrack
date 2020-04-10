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
                <div>
                  <b-tag type="is-info" v-if="list.completed">Completed</b-tag>
                  <b-tag type="is-warning" v-if="!list.completed">Uncompleted</b-tag>
                </div>
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
          </div>
          <div class="content" v-if="list.notes">
            {{ list.notes }}
            <br/>
            <br/>
            <div v-if="typeof list.count !== 'undefined' && list.count > 0">
              {{ list.count }} item(s)
            </div>
            <div v-else>
              0 items
            </div>
          </div>
        </div>
      </div>
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
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    }
  },
  async beforeMount () {
    var list = this.list
    var userId = list.author

    flatmates.GetFlatmate(userId).then(resp => {
      this.authorNames = resp.data.spec.names
      return flatmates.GetFlatmate(this.list.authorLast)
    }).then(resp => {
      this.authorLastNames = resp.data.spec.names
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
