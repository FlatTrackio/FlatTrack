<template>
  <div>
    <div v-for="list in lists" v-bind:key="list">
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
                <p class="title is-4">{{ list.name }}</p>
                <b-tag type="is-info" v-if="list.completed">Completed</b-tag>
                <b-tag type="is-warning" v-if="!list.completed">Uncompleted</b-tag>
                <p class="subtitle is-6">
                  <span v-if="list.creationTimestamp == list.modificationTimestamp">
                    Created
                  </span>
                  <span v-else>
                    Updated
                  </span>
                  {{ TimestampToCalendar(list.creationTimestamp) }}, by {{ list.author }}
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
  </div>
</template>

<script>
import common from '@/frontend/common/common'

export default {
  name: 'shopping-list-card-view',
  props: {
    lists: Object,
    authors: Object
  },
  methods: {
    goToRef (ref) {
      this.$router.push({ path: ref })
    },
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    }
  }
}
</script>
