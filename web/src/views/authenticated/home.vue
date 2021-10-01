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
    <div class="container">
      <section class="section">
        <b-loading :is-full-page="false" :active.sync="pageLoading" :can-cancel="false"></b-loading>
        <h1 class="title is-1">Home</h1>
        <p class="subtitle is-3">Recent activity</p>
        <h1 class="subtitle is-4">Shopping lists</h1>
        <div v-if="lists.length > 0">
          <shoppingListCardView :v-motion-slide-bottom="enableAnimations" :list="list" :authors="authors" :lists="lists" :index="index" v-for="(list, index) in lists" v-bind:key="list" :deviceIsMobile="true" :mini="true" />
          <p class="is-size-5 ml-3 mb-2">
            <b-icon icon="party-popper" type="is-success" size="is-medium"></b-icon>
            You're all caught up! That's all for now
          </p>
          <b-button
            class="has-text-left"
            @click="ClearItems"
            type="is-info"
            icon-left="close-box-outline"
            expanded>
            Clear
          </b-button>
        </div>
        <div v-else>
          <b-message type="is-warning">
            You're all caught up. Check back later!
          </b-message>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/common/common'
import shoppinglist from '@/requests/authenticated/shoppinglist'
import dayjs from 'dayjs'

export default {
  name: 'home',
  data () {
    return {
      enableAnimations: common.GetEnableAnimations() === 'true',
      authors: {},
      pageLoading: true,
      lists: [],
      deviceIsMobile: false,
      modificationTimestampAfter: common.GetHomeLastViewedTimestamp(),
      sortBy: 'recentlyUpdated'
    }
  },
  methods: {
    GetShoppingLists () {
      shoppinglist.GetShoppingLists(undefined, this.sortBy, undefined, this.modificationTimestampAfter, 5).then(resp => {
        this.pageLoading = false
        this.lists = resp.data.list || []
      }).catch(() => {
        common.DisplayFailureToast('Hmmm seems somethings gone wrong loading the shopping lists')
      })
    },
    ClearItems () {
      common.WriteHomeLastViewedTimestamp(Number(dayjs().unix()))
      this.lists = []
    }
  },
  components: {
    shoppingListCardView: () => import('@/components/authenticated/shopping-list-card-view.vue')
  },
  async beforeMount () {
    this.GetShoppingLists()
  }
}
</script>

<style scoped>
</style>
