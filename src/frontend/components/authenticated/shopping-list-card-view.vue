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
              <div class="block">
                <p class="title is-4">{{ list.name }}</p>
                <b-tag type="is-info" v-if="list.completed">Completed</b-tag>
                <b-tag type="is-warning" v-if="!list.completed">Uncompleted</b-tag>
              </div>
              <p class="subtitle is-6">
                <span v-if="list.creationTimestamp == list.modificationTimestamp">
                  Created
                </span>
                <span v-else>
                  Updated
                </span>
                {{ TimestampToCalendar(list.creationTimestamp) }}, by {{ authorNames }}
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
      authorNames: ''
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
  async created () {
    var list = this.list
    var userId = list.author
    flatmates.GetFlatmate(userId).then(resp => {
      this.authorNames = resp.data.spec.names
    }).catch(err => {
      common.DisplayFailureToast('Unable to find author of list' + '<br/>' + err.response.data.metadata.response)
    })
  }
}
</script>
